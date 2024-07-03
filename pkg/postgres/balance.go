package postgres

import (
	"context"

	"lystem/internal/models/balance"
	"lystem/internal/models/user"
	"lystem/internal/models/withdrawal"
	"lystem/internal/repository"
)

func (s *DBStorage) FindBalance(ctx context.Context, currentUser *user.User) (*balance.Balance, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	balancesRepo := repository.NewBalancesRepository(conn)
	userBalance, err := balancesRepo.FindByUser(ctx, currentUser)
	if err != nil {
		return nil, newDBError(err)
	}
	return userBalance, nil
}

//func (s *DBStorage) IncreaseBalance(ctx context.Context, processedOrder *order.Order) error {
//	conn, err := s.instance.Acquire(ctx)
//	if err != nil {
//		return newDBError(err)
//	}
//	defer conn.Release()
//	balancesRepo := repository.NewBalancesRepository(conn)
//
//	if err = balancesRepo.Accrual(ctx, processedOrder); err != nil {
//		return newDBError(err)
//	}
//	return nil
//}

func (s *DBStorage) DeductFromBalance(ctx context.Context, w *withdrawal.Withdrawal, currUser *user.User) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return newDBError(err)
	}
	defer conn.Release()
	balancesRepo := repository.NewBalancesRepository(conn)

	err = balancesRepo.Decrease(ctx, w, currUser)
	if err != nil {
		return newDBError(err)
	}
	return nil
}

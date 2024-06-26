package postgres

import (
	"context"

	"lystem/internal/models/balance"
	"lystem/internal/models/order"
	"lystem/internal/models/user"
	"lystem/internal/models/withdrawal"
	"lystem/internal/repository"
)

func (s *DBStorage) IncreaseBalance(ctx context.Context, processedOrder *order.Order) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return newDBError(err)
	}
	defer conn.Release()
	balancesRepo := repository.NewBalancesRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return newDBError(err)
	}

	err = balancesRepo.Accrual(ctx, tx, processedOrder)
	if err != nil {
		return rollbackOnErr(ctx, tx, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return newDBError(err)
	}
	return nil
}

func (s *DBStorage) DeductFromBalance(ctx context.Context, w *withdrawal.Withdrawal, currUser *user.User) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return newDBError(err)
	}
	defer conn.Release()
	balancesRepo := repository.NewBalancesRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return newDBError(err)
	}

	err = balancesRepo.Decrease(ctx, tx, w, currUser)
	if err != nil {
		return rollbackOnErr(ctx, tx, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return newDBError(err)
	}
	return nil
}
func (s *DBStorage) FindBalance(ctx context.Context, currentUser *user.User) (*balance.Balance, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	balancesRepo := repository.NewBalancesRepository(conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	userBalance, err := balancesRepo.FindByUser(ctx, tx, currentUser)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return userBalance, nil
}

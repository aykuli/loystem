package postgres

import (
	"context"

	"lystem/internal/models/balance"
	"lystem/internal/models/user"
	"lystem/internal/models/withdrawal"
	"lystem/internal/repository"
)

func (s *DBStorage) FindWithdrawals(ctx context.Context, userBalance *balance.Balance) ([]withdrawal.Withdrawal, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	withdrawalsRepo := repository.NewWithdrawalsRepository(conn)
	withdrawals, err := withdrawalsRepo.FindAll(ctx, userBalance)
	if err != nil {
		return nil, newDBError(err)
	}
	return withdrawals, nil
}

func (s *DBStorage) CreateWithdrawal(ctx context.Context, orderNumber string, currUser *user.User, sum float64) (*withdrawal.Withdrawal, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	withdrawalsRepo := repository.NewWithdrawalsRepository(conn)
	balancesRepo := repository.NewBalancesRepository(conn)

	userBalance, err := balancesRepo.FindByUser(ctx, currUser)
	if err != nil {
		return nil, newDBError(err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	withdraw, err := withdrawalsRepo.Create(ctx, tx, orderNumber, userBalance, sum)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}
	if err = balancesRepo.Decrease(ctx, withdraw, currUser); err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return withdraw, nil
}

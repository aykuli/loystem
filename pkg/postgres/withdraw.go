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
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	withdrawals, err := withdrawalsRepo.FindAll(ctx, tx, userBalance)
	if err != nil {
		return nil, newDBError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return withdrawals, nil
}

func (s *DBStorage) CreateWithdraw(ctx context.Context, orderNumber string, currUser *user.User, sum float64) (*withdrawal.Withdrawal, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	withdrawalsRepo := repository.NewWithdrawalsRepository(conn)
	balanceRepo := repository.NewBalancesRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}

	userBalance, err := balanceRepo.FindByUser(ctx, tx, currUser)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}

	withdraw, err := withdrawalsRepo.Create(ctx, tx, orderNumber, userBalance, sum)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}

	if err = balanceRepo.Decrease(ctx, tx, withdraw, currUser); err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return withdraw, nil
}

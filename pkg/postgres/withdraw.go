package postgres

import (
	"github.com/valyala/fasthttp"

	"lystem/internal/models/balance"
	"lystem/internal/models/order"
	"lystem/internal/models/user"
	"lystem/internal/models/withdrawal"
	"lystem/internal/repository"
)

func (s *DBStorage) FindWithdrawals(ctx *fasthttp.RequestCtx, userBalance *balance.Balance) ([]withdrawal.Withdrawal, error) {
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

func (s *DBStorage) CreateWithdraw(ctx *fasthttp.RequestCtx, processedOrder *order.Order, currUser *user.User, sum float64) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return newDBError(err)
	}
	defer conn.Release()

	withdrawalsRepo := repository.NewWithdrawalsRepository(conn)
	balanceRepo := repository.NewBalancesRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return newDBError(err)
	}

	userBalance, err := balanceRepo.FindByUser(ctx, tx, currUser)
	if err != nil {
		return rollbackOnErr(ctx, tx, err)
	}

	withdraw, err := withdrawalsRepo.Create(ctx, tx, processedOrder, userBalance, sum)
	if err != nil {
		return rollbackOnErr(ctx, tx, err)
	}

	err = balanceRepo.Decrease(ctx, tx, withdraw, currUser)
	if err != nil {
		return rollbackOnErr(ctx, tx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return newDBError(err)
	}
	return nil
}

package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/valyala/fasthttp"

	"lystem/internal/models/order"
	"lystem/internal/models/user"
	"lystem/internal/repository"
)

func (s *DBStorage) FindOrderByNumber(ctx *fasthttp.RequestCtx, number string) (*order.Order, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	ordersRepo := repository.NewOrdersRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	foundOrder, err := ordersRepo.FindByNumber(ctx, tx, number)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return nil, newDBError(rollbackErr)
		}
		return nil, nil
	}

	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}

	return foundOrder, nil
}

func (s *DBStorage) SaveOrder(ctx *fasthttp.RequestCtx, number string, userID int) (*order.Order, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	ordersRepo := repository.NewOrdersRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	savedOrder, err := ordersRepo.Save(ctx, tx, number, userID)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return savedOrder, nil
}

func (s *DBStorage) UpdateOrder(ctx *fasthttp.RequestCtx, newOrder *order.Order) (*order.Order, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	ordersRepo := repository.NewOrdersRepository(conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	updatedOrder, err := ordersRepo.Update(ctx, tx, newOrder)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return updatedOrder, nil
}

func (s *DBStorage) SelectOrders(ctx *fasthttp.RequestCtx, currentUser *user.User) ([]order.Order, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()
	ordersRepo := repository.NewOrdersRepository(conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, newDBError(err)
	}

	orders, err := ordersRepo.SelectAll(ctx, tx, currentUser)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return orders, nil
}

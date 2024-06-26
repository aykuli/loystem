package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"lystem/internal/models/order"
	"lystem/internal/models/user"
	"lystem/internal/repository"
)

func (s *DBStorage) FindOrderByNumber(ctx context.Context, number string) (*order.Order, error) {
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

func (s *DBStorage) SaveOrder(ctx context.Context, number string, userID int) (*order.Order, error) {
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

func (s *DBStorage) UpdateOrder(ctx context.Context, newOrder *order.Order) (*order.Order, error) {
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

func (s *DBStorage) SelectUserOrders(ctx context.Context, u *user.User) ([]order.Order, error) {
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

	orders, err := ordersRepo.SelectUserOrders(ctx, tx, u)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return orders, nil
}

func (s *DBStorage) SelectAccrualOrders(ctx context.Context) ([]order.Order, error) {
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

	orders, err := ordersRepo.SelectAccrualOrders(ctx, tx)
	if err != nil {
		return nil, rollbackOnErr(ctx, tx, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, newDBError(err)
	}
	return orders, nil
}

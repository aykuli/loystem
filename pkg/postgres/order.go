package postgres

import (
	"context"
	"errors"
	"fmt"

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
	fmt.Println("1")

	ordersRepo := repository.NewOrdersRepository(conn)

	foundOrder, err := ordersRepo.FindByNumber(ctx, number)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
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

	savedOrder, err := ordersRepo.Save(ctx, number, userID)
	if err != nil {
		return nil, newDBError(err)
	}
	return savedOrder, nil
}

func (s *DBStorage) UpdateOrder(ctx context.Context, newOrder *order.Order) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return newDBError(err)
	}
	defer conn.Release()

	ordersRepo := repository.NewOrdersRepository(conn)
	err = ordersRepo.Update(ctx, newOrder)
	if err != nil {
		return newDBError(err)
	}
	return nil
}

func (s *DBStorage) UpdateOrderAndIncreaseBalance(ctx context.Context, newOrder *order.Order) error {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return newDBError(err)
	}
	defer conn.Release()

	ordersRepo := repository.NewOrdersRepository(conn)
	balancesRepo := repository.NewBalancesRepository(conn)

	tx, err := conn.Begin(ctx)
	if err != nil {
		return newDBError(err)
	}

	updatedOrder, err := ordersRepo.UpdateReturning(ctx, tx, newOrder)
	if err != nil {
		return rollbackOnErr(ctx, tx, err)
	}
	if err = balancesRepo.Accrual(ctx, tx, updatedOrder); err != nil {
		return rollbackOnErr(ctx, tx, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return newDBError(err)
	}

	return nil
}

func (s *DBStorage) FindAllUserOrders(ctx context.Context, u *user.User) ([]order.Order, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	ordersRepo := repository.NewOrdersRepository(conn)
	orders, err := ordersRepo.FindAllUserOrders(ctx, u)
	if err != nil {
		return nil, newDBError(err)
	}

	return orders, nil
}

func (s *DBStorage) SelectUnprocessedOrders(ctx context.Context, limit int) ([]order.Order, error) {
	conn, err := s.instance.Acquire(ctx)
	if err != nil {
		return nil, newDBError(err)
	}
	defer conn.Release()

	ordersRepo := repository.NewOrdersRepository(conn)
	orders, err := ordersRepo.SelectUnprocessed(ctx, limit)
	if err != nil {
		return nil, newDBError(err)
	}
	return orders, nil
}

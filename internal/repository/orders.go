package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lystem/internal/models/order"
	"lystem/internal/models/user"
)

var (
	selectOrderByNumberSQL  = `SELECT id, number, user_id, status FROM orders WHERE number = @number`
	insertOrderSQL          = `INSERT INTO orders (number, user_id, accrual, status) VALUES (@number, @user_id, @accrual, @status) RETURNING id`
	updateOrderSQL          = `UPDATE orders SET (accrual, status) = (@accrual, @status) WHERE number = @number RETURNING id, number, user_id, status, accrual, uploaded_at`
	selectOrdersByUserIDSQL = `SELECT id, number, user_id, status,accrual FROM orders WHERE user_id = @user_id`
)

type OrdersRepository struct {
	conn *pgxpool.Conn
}

func NewOrdersRepository(conn *pgxpool.Conn) *OrdersRepository {
	return &OrdersRepository{conn: conn}
}

func (r *OrdersRepository) FindByNumber(ctx context.Context, tx pgx.Tx, number string) (*order.Order, error) {
	args := pgx.NamedArgs{"number": number}
	result := tx.QueryRow(ctx, selectOrderByNumberSQL, args)
	var scannedOrder order.Order
	if err := result.Scan(&scannedOrder.ID, &scannedOrder.Number, &scannedOrder.UserID, &scannedOrder.Status); err != nil {
		return nil, err
	}
	return &scannedOrder, nil
}

func (r *OrdersRepository) Save(ctx context.Context, tx pgx.Tx, number string, userID int) (*order.Order, error) {
	args := pgx.NamedArgs{"number": number, "user_id": userID, "accrual": 0, "status": order.StatusNew}
	result := tx.QueryRow(ctx, insertOrderSQL, args)
	var id int
	if err := result.Scan(&id); err != nil {
		return nil, err
	}

	return &order.Order{ID: id, Number: number, UserID: userID}, nil
}

func (r *OrdersRepository) Update(ctx context.Context, tx pgx.Tx, newOrder *order.Order) (*order.Order, error) {
	args := pgx.NamedArgs{"number": newOrder.Number, "accrual": newOrder.Accrual, "status": newOrder.Status}
	result := tx.QueryRow(ctx, updateOrderSQL, args)
	var scannedOrder order.Order
	if err := result.Scan(&scannedOrder.ID, &scannedOrder.Number, &scannedOrder.UserID, &scannedOrder.Status, &scannedOrder.Accrual, &scannedOrder.UploadedAt); err != nil {
		return nil, err
	}
	return &scannedOrder, nil
}

func (r *OrdersRepository) SelectAll(ctx context.Context, tx pgx.Tx, currentUser *user.User) ([]order.Order, error) {
	rows, err := tx.Query(ctx, selectOrdersByUserIDSQL, pgx.NamedArgs{"user_id": currentUser.ID})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []order.Order
	for rows.Next() {
		var sOrder order.Order
		if err = rows.Scan(&sOrder.ID, &sOrder.Number, &sOrder.UserID, &sOrder.Status, &sOrder.Accrual); err != nil {
			return nil, err
		}
		orders = append(orders, sOrder)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

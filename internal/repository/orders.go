package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lystem/internal/models/order"
	"lystem/internal/models/user"
)

var (
	selectOrderByNumberSQL    = `SELECT id, number, user_id, status FROM orders WHERE number = @number`
	insertOrderSQL            = `INSERT INTO orders (number, user_id, accrual, status) VALUES (@number, @user_id, @accrual, @status) RETURNING id`
	updateOrderSQL            = `UPDATE orders SET (accrual, status) = (@accrual, @status) WHERE number = @number`
	updateReturningOrderSQL   = `UPDATE orders SET (accrual, status) = (@accrual, @status) WHERE number = @number RETURNING id, number, user_id, status, accrual, uploaded_at`
	selectOrdersByUserIDSQL   = `SELECT id, number, user_id, status, accrual, uploaded_at FROM orders WHERE user_id = @user_id`
	selectOrdersByStatusesSQL = `SELECT id, number, user_id, status, accrual, uploaded_at FROM orders WHERE status IN ('NEW','REGISTERED','PROCESSING') LIMIT @limit`
)

type OrdersRepository struct {
	conn *pgxpool.Conn
}

func NewOrdersRepository(conn *pgxpool.Conn) *OrdersRepository {
	return &OrdersRepository{conn: conn}
}

func (r *OrdersRepository) FindByNumber(ctx context.Context, number string) (*order.Order, error) {
	result := r.conn.QueryRow(ctx, selectOrderByNumberSQL, pgx.NamedArgs{"number": number})
	var scannedOrder order.Order
	if err := result.Scan(&scannedOrder.ID, &scannedOrder.Number, &scannedOrder.UserID, &scannedOrder.Status); err != nil {
		return nil, err
	}
	return &scannedOrder, nil
}

func (r *OrdersRepository) Save(ctx context.Context, number string, userID int) (*order.Order, error) {
	args := pgx.NamedArgs{"number": number, "user_id": userID, "accrual": 0, "status": order.StatusNew}
	result := r.conn.QueryRow(ctx, insertOrderSQL, args)
	var id int
	if err := result.Scan(&id); err != nil {
		return nil, err
	}

	return &order.Order{ID: id, Number: number, UserID: userID}, nil
}

func (r *OrdersRepository) Update(ctx context.Context, newOrder *order.Order) error {
	args := pgx.NamedArgs{"number": newOrder.Number, "accrual": newOrder.Accrual, "status": newOrder.Status}
	_, err := r.conn.Exec(ctx, updateOrderSQL, args)
	return err
}

func (r *OrdersRepository) UpdateReturning(ctx context.Context, tx pgx.Tx, newOrder *order.Order) (*order.Order, error) {
	args := pgx.NamedArgs{"number": newOrder.Number, "accrual": newOrder.Accrual, "status": newOrder.Status}
	result := tx.QueryRow(ctx, updateReturningOrderSQL, args)
	var scannedOrder order.Order
	if err := result.Scan(&scannedOrder.ID, &scannedOrder.Number, &scannedOrder.UserID, &scannedOrder.Status, &scannedOrder.Accrual, &scannedOrder.UploadedAt); err != nil {
		return nil, err
	}

	return &scannedOrder, nil
}

func (r *OrdersRepository) FindAllUserOrders(ctx context.Context, u *user.User) ([]order.Order, error) {
	rows, err := r.conn.Query(ctx, selectOrdersByUserIDSQL, pgx.NamedArgs{"user_id": u.ID})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []order.Order
	for rows.Next() {
		var sOrder order.Order
		if err = rows.Scan(&sOrder.ID, &sOrder.Number, &sOrder.UserID, &sOrder.Status, &sOrder.Accrual, &sOrder.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, sOrder)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrdersRepository) SelectUnprocessed(ctx context.Context, limit int) ([]order.Order, error) {
	rows, err := r.conn.Query(ctx, selectOrdersByStatusesSQL, pgx.NamedArgs{"limit": limit})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []order.Order
	for rows.Next() {
		var sOrder order.Order
		if err = rows.Scan(&sOrder.ID, &sOrder.Number, &sOrder.UserID, &sOrder.Status, &sOrder.Accrual, &sOrder.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, sOrder)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

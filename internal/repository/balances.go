package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lystem/internal/models/balance"
	"lystem/internal/models/order"
	"lystem/internal/models/user"
	"lystem/internal/models/withdrawal"
)

var (
	insertBalanceSQL     = `INSERT INTO balances (user_id, current) VALUES (@user_id, @current)`
	selectBalanceSQL     = `SELECT current, user_id FROM balances WHERE user_id = @user_id`
	increaseBalanceSQL   = `UPDATE balances SET current = current + @accrual WHERE user_id = @user_id`
	deductFromBalanceSQL = `UPDATE balances SET current = current - @sum WHERE user_id = @user_id`
)

type BalancesRepository struct {
	conn *pgxpool.Conn
}

func NewBalancesRepository(conn *pgxpool.Conn) *BalancesRepository {
	return &BalancesRepository{conn}
}

func (r *BalancesRepository) Create(ctx context.Context, tx pgx.Tx, currUser *user.User) (*balance.Balance, error) {
	args := pgx.NamedArgs{"user_id": currUser.ID, "current": 0}
	if _, err := tx.Exec(ctx, insertBalanceSQL, args); err != nil {
		return nil, err
	}
	return &balance.Balance{Current: 0, UserID: currUser.ID}, nil
}

func (r *BalancesRepository) FindByUser(ctx context.Context, currUser *user.User) (*balance.Balance, error) {
	args := pgx.NamedArgs{"user_id": currUser.ID}
	result := r.conn.QueryRow(ctx, selectBalanceSQL, args)
	var b balance.Balance
	if err := result.Scan(&b.Current, &b.UserID); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BalancesRepository) Accrual(ctx context.Context, tx pgx.Tx, o *order.Order) error {
	args := pgx.NamedArgs{"accrual": o.Accrual, "user_id": o.UserID}
	if _, err := tx.Exec(ctx, increaseBalanceSQL, args); err != nil {
		return err
	}
	return nil
}

func (r *BalancesRepository) Decrease(ctx context.Context, w *withdrawal.Withdrawal, currUser *user.User) error {
	args := pgx.NamedArgs{"sum": w.Sum, "user_id": currUser.ID}
	if _, err := r.conn.Exec(ctx, deductFromBalanceSQL, args); err != nil {
		return err
	}
	return nil
}

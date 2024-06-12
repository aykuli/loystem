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
	insertBalanceSQL     = `INSERT INTO balances (user_id, 0) RETURNING id`
	selectBalanceSQL     = `SELECT id, current, user_id FROM balances WHERE user_id = @user_id`
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
	args := pgx.NamedArgs{"user_id": currUser.ID}
	result := tx.QueryRow(ctx, insertBalanceSQL, args)
	var id int
	if err := result.Scan(&id); err != nil {
		return nil, err
	}
	return &balance.Balance{ID: id, Current: 0, UserID: currUser.ID}, nil
}

func (r *BalancesRepository) FindByUser(ctx context.Context, tx pgx.Tx, currUser *user.User) (*balance.Balance, error) {
	args := pgx.NamedArgs{"user_id": currUser.ID}
	result := tx.QueryRow(ctx, selectBalanceSQL, args)
	var b balance.Balance
	if err := result.Scan(&b.ID, &b.Current, &b.UserID); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BalancesRepository) Increase(ctx context.Context, tx pgx.Tx, o *order.Order, currUser *user.User) error {
	args := pgx.NamedArgs{"accrual": o.Accrual, "user_id": currUser.ID}
	if _, err := tx.Exec(ctx, increaseBalanceSQL, args); err != nil {
		return err
	}
	return nil
}

func (r *BalancesRepository) Decrease(ctx context.Context, tx pgx.Tx, w *withdrawal.Withdrawal, currUser *user.User) error {
	args := pgx.NamedArgs{"sum": w.Sum, "user_id": currUser.ID}
	if _, err := tx.Exec(ctx, deductFromBalanceSQL, args); err != nil {
		return err
	}
	return nil
}

package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lystem/internal/models/balance"
	"lystem/internal/models/withdrawal"
)

var (
	insertWithdrawalSQL  = `INSERT INTO withdrawals (sum, order_number, balance_id) VALUES (@sum, @order_number, @balance_id) RETURNING id, sum, balance_id`
	selectWithdrawalsSQL = `SELECT id, sum, order_number, balance_id FROM withdrawals WHERE balance_id = @balance_id`
)

type WithdrawalsRepository struct {
	conn *pgxpool.Conn
}

func NewWithdrawalsRepository(conn *pgxpool.Conn) *WithdrawalsRepository {
	return &WithdrawalsRepository{conn}
}

func (r *WithdrawalsRepository) FindAll(ctx context.Context, b *balance.Balance) ([]withdrawal.Withdrawal, error) {
	rows, err := r.conn.Query(ctx, selectWithdrawalsSQL, pgx.NamedArgs{"balance_id": b.UserID})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []withdrawal.Withdrawal
	for rows.Next() {
		var w withdrawal.Withdrawal
		if err = rows.Scan(&w.ID, &w.Sum, &w.OrderNumber, &w.BalanceID); err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (r *WithdrawalsRepository) Create(ctx context.Context, tx pgx.Tx, orderNumber string, userBalance *balance.Balance, sum float64) (*withdrawal.Withdrawal, error) {
	args := pgx.NamedArgs{"order_number": orderNumber, "balance_id": userBalance.UserID, "sum": sum}
	result := tx.QueryRow(ctx, insertWithdrawalSQL, args)
	var wd = withdrawal.Withdrawal{OrderNumber: orderNumber}
	if err := result.Scan(&wd.ID, &wd.Sum, &wd.BalanceID); err != nil {
		return nil, err
	}

	return &wd, nil
}

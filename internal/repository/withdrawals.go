package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lystem/internal/models/balance"
	"lystem/internal/models/order"
	"lystem/internal/models/withdrawal"
)

var (
	insertWithdrawalSQL  = `INSERT INTO withdrawals (sum, order_id, balance_id) VALUES (@sum, @order_id, @balance_id) RETURNING id, sum, order_id, balance_id`
	selectWithdrawalsSQL = `SELECT id, sum, order_id, balance_id FROM withdrawals WHERE balance_id = @balance_id`
)

type WithdrawalsRepository struct {
	conn *pgxpool.Conn
}

func NewWithdrawalsRepository(conn *pgxpool.Conn) *WithdrawalsRepository {
	return &WithdrawalsRepository{conn}
}

func (r *WithdrawalsRepository) FindAll(ctx context.Context, tx pgx.Tx, b *balance.Balance) ([]withdrawal.Withdrawal, error) {
	rows, err := tx.Query(ctx, selectWithdrawalsSQL, pgx.NamedArgs{"balance_id": b.ID})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []withdrawal.Withdrawal
	for rows.Next() {
		var w withdrawal.Withdrawal
		if err = rows.Scan(&w.ID, &w.Sum, &w.OrderID, &w.BalanceID); err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return withdrawals, nil
}

func (r *WithdrawalsRepository) Create(ctx context.Context, tx pgx.Tx, processedOrder *order.Order, userBalance *balance.Balance, sum float64) (*withdrawal.Withdrawal, error) {
	args := pgx.NamedArgs{"order_id": processedOrder.ID, "balance_id": userBalance.ID, "sum": sum}
	result := tx.QueryRow(ctx, insertWithdrawalSQL, args)
	var wd withdrawal.Withdrawal
	if err := result.Scan(&wd.ID, &wd.Sum, &wd.OrderID, &wd.BalanceID); err != nil {
		return nil, err
	}

	return &wd, nil
}

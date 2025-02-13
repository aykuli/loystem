package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	createUsersTableSQL = `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login TEXT NOT NULL,
		hashed_password TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT now())`
	createUsersLoginKeySQL = `CREATE UNIQUE INDEX IF NOT EXISTS users_login_key ON users(login)`
	createSessionsTableSQL = `CREATE TABLE IF NOT EXISTS sessions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now())`
	createOrdersTableSQL = `CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		number VARCHAR NOT NULL UNIQUE,
		accrual FLOAT,
		status VARCHAR DEFAULT 'NEW',
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		UNIQUE (number, user_id)
	)`
	createTableBalancesSQL = `CREATE TABLE IF NOT EXISTS balances (
		user_id INTEGER PRIMARY KEY REFERENCES users(id),
		current FLOAT NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`
	createTableWithdrawalsSQL = `CREATE TABLE IF NOT EXISTS withdrawals (
		id SERIAL PRIMARY KEY,
		sum FLOAT NOT NULL,
		balance_id INTEGER NOT NULL REFERENCES balances(user_id),
		order_number VARCHAR NOT NULL,
		proceeded_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`
)

type Repository struct {
	conn *pgxpool.Conn
}

func New(conn *pgxpool.Conn) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) Init(ctx context.Context, tx pgx.Tx) error {
	queries := []string{
		createUsersTableSQL,
		createUsersLoginKeySQL,
		createSessionsTableSQL,
		createOrdersTableSQL,
		createTableBalancesSQL,
		createTableWithdrawalsSQL,
	}
	for _, query := range queries {
		if _, err := tx.Exec(ctx, query); err != nil {
			return err
		}
	}
	return nil
}

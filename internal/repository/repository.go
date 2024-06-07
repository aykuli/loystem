package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	createUsersTableSQL = `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login TEXT NOT NULL,
		salt TEXT NOT NULL,
		hashed_password TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT now())`
	createUsersLoginKeySQL = `CREATE UNIQUE INDEX IF NOT EXISTS users_login_key ON users(login)`
	createSessionsTableSQL = `CREATE TABLE IF NOT EXISTS sessions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now())`
	createOrderStatusTypeSQL = `CREATE TYPE status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')`
	createOrdersTableSQL     = `CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		number VARCHAR NOT NULL UNIQUE,
		accrual FLOAT,
		status status,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		UNIQUE (number, user_id)
	)`
	createTableBalancesSQL = `CREATE TABLE IF NOT EXISTS balances (
		id SERIAL PRIMARY KEY,
		current FLOAT NOT NULL,
		user_id INTEGER NOT NULL REFERENCES users(id),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`
	createBalanceOperationTypeSQL = `CREATE TYPE operation AS ENUM('withdrawn', 'earned', 'summarized')`
	createTableBalanceEventsSQL   = `CREATE TABLE IF NOT EXISTS withdrawals (
		id SERIAL PRIMARY KEY,
		sum FLOAT NOT NULL,
		operation operation,
		balance_id INTEGER NOT NULL REFERENCES balances(id),
		order_id INTEGER NOT NULL REFERENCES orders(id),
		proceeded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		UNIQUE (order_id, balance_id)
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
		createOrderStatusTypeSQL,
		createOrdersTableSQL,
		createTableBalancesSQL,
		createBalanceOperationTypeSQL,
		createTableBalanceEventsSQL,
	}
	for _, query := range queries {
		if _, err := tx.Exec(ctx, query); err != nil {
			fmt.Println("\n\n", query)
			fmt.Println(err)
			return err
		}
	}
	return nil
}

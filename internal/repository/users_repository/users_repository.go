package users_repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lystem/internal/model/user"
)

var (
	createUsersTableSQL = `CREATE TABLE IF NOT EXISTS users ( id SERIAL PRIMARY KEY,
 login TEXT NOT NULL,
 salt TEXT NOT NULL,
 hashed_password TEXT NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT now()) `
	insertUserSQL      = "INSERT INTO users (login, salt, hashed_password) VALUES (@login, @salt, @hashed_password) RETURNING id, login, salt, hashed_password"
	findUserByLoginSQL = "SELECT id, salt, hashed_password FROM users WHERE login=@login"
)

type Repository struct {
	conn *pgxpool.Conn
}

func New(conn *pgxpool.Conn) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) Create(salt, login, hashedPassword string) (*user.User, error) {
	var u user.User
	args := pgx.NamedArgs{"login": login, "salt": salt, "hashed_password": hashedPassword}
	result := r.conn.QueryRow(context.Background(), insertUserSQL, args)
	err := result.Scan(&u.ID, &u.Login, &u.Salt, &u.HashedPassword)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FindByLogin(login string) (*user.User, error) {
	var u user.User
	args := pgx.NamedArgs{"login": login}
	result := r.conn.QueryRow(context.Background(), findUserByLoginSQL, args)
	err := result.Scan(&u.ID, &u.Login, &u.Salt, &u.HashedPassword)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

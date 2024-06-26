package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lystem/internal/models/user"
)

var (
	insertUserSQL      = `INSERT INTO users (login, salt, hashed_password) VALUES (@login, @salt, @hashed_password) RETURNING id`
	findUserByLoginSQL = `SELECT id, login, hashed_password, salt FROM users WHERE login = @login`
	findUserByIDSQL    = `SELECT id FROM users WHERE id = @id`
)

type UsersRepository struct {
	conn *pgxpool.Conn
}

func NewUsersRepository(conn *pgxpool.Conn) *UsersRepository {
	return &UsersRepository{conn: conn}
}

func (r *UsersRepository) Create(ctx context.Context, tx pgx.Tx, u *user.User) (*user.User, error) {
	args := pgx.NamedArgs{"login": u.Login, "salt": u.Salt, "hashed_password": u.HashedPassword}
	result := tx.QueryRow(ctx, insertUserSQL, args)

	var id int
	if err := result.Scan(&id); err != nil {
		return nil, err
	}
	return &user.User{
		ID:             id,
		Salt:           u.Salt,
		Login:          u.Login,
		HashedPassword: u.HashedPassword,
	}, nil
}

func (r *UsersRepository) FindByLogin(login string) (*user.User, error) {
	var u user.User
	result := r.conn.QueryRow(context.Background(), findUserByLoginSQL, pgx.NamedArgs{"login": login})
	if err := result.Scan(&u.ID, &u.Login, &u.HashedPassword, &u.Salt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepository) FindByID(id int) (*user.User, error) {
	var u user.User
	result := r.conn.QueryRow(context.Background(), findUserByIDSQL, pgx.NamedArgs{"id": id})
	if err := result.Scan(&u.ID); err != nil {
		return nil, err
	}
	return &u, nil
}

package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lystem/internal/models/user"
)

var (
	insertUserSQL      = "INSERT INTO users (login, salt, hashed_password) VALUES (@login, @salt, @hashed_password) RETURNING id"
	findUserByLoginSQL = "SELECT id, login, hashed_password, salt FROM users WHERE login = @login"
	findUserByIDSQL    = "SELECT id FROM users WHERE id = @id"
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

	var id int64
	err := result.Scan(&id)
	if err != nil {
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
	args := pgx.NamedArgs{"login": login}
	result := r.conn.QueryRow(context.Background(), findUserByLoginSQL, args)
	err := result.Scan(&u.ID, &u.Login, &u.HashedPassword, &u.Salt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepository) FindByID(id int64) (*user.User, error) {
	var u user.User
	args := pgx.NamedArgs{"id": id}
	result := r.conn.QueryRow(context.Background(), findUserByIDSQL, args)
	err := result.Scan(&u.ID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lystem/internal/models/session"
	"lystem/internal/models/user"
)

var (
	insertSQL   = "INSERT INTO sessions (user_id) VALUES (@user_id) RETURNING id, created_at"
	deleteSQL   = "DELETE FROM sessions WHERE user_id = @user_id"
	findByIDSQL = "SELECT id, user_id, created_at FROM sessions WHERE user_id = @user_id"
)

type SessionsRepository struct {
	conn *pgxpool.Conn
}

func NewSessionsRepository(conn *pgxpool.Conn) *SessionsRepository {
	return &SessionsRepository{conn: conn}
}

func (r *SessionsRepository) Create(ctx context.Context, tx pgx.Tx, u *user.User) (*session.Session, error) {
	result := tx.QueryRow(ctx, insertSQL, pgx.NamedArgs{"user_id": u.ID})
	var newSession session.Session
	err := result.Scan(&newSession.ID, &newSession.CreatedAt)
	if err != nil {
		return nil, err
	}

	// returns new sessions and deletes all user's prev sessions
	return &newSession, r.Delete(ctx, tx, u)
}

func (r *SessionsRepository) FindByID(ctx context.Context, tx pgx.Tx, id string) (*session.Session, error) {
	result := tx.QueryRow(ctx, findByIDSQL, pgx.NamedArgs{"id": id})
	var foundSession session.Session
	err := result.Scan(&foundSession.ID, &foundSession.UserID, &foundSession.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &foundSession, nil
}

func (r *SessionsRepository) Delete(ctx context.Context, tx pgx.Tx, u *user.User) error {
	_, err := tx.Exec(ctx, deleteSQL, pgx.NamedArgs{"user_id": u.ID})
	return err
}
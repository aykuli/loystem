package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type postgresError struct {
	name string
	Err  error
}

var (
	ErrUserAlreadyExists = errors.New("логин уже занят")
)

func (dbErr *postgresError) Error() string {
	return fmt.Sprintf("[%s] %v", dbErr.name, dbErr.Err.Error())
}

func newDBError(err error) error {
	var pgErr *pgconn.PgError
	name := "Postgres"
	if errors.As(err, &pgErr) {
		if pgerrcode.IsConnectionException(pgErr.Code) {
			name = "Postgres. Connection Exception"
		}

		if pgerrcode.IsDataException(pgErr.Code) {
			name = "Postgres. Data Exception"
		}

	}

	return &postgresError{
		name: name,
		Err:  err,
	}
}

func rollbackOnErr(ctx context.Context, tx pgx.Tx, err error) error {
	if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
		return newDBError(rollbackErr)
	}

	return newDBError(err)
}

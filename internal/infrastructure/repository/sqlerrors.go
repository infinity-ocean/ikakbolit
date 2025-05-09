package repository

import (
	"errors"
	"fmt"

	"github.com/infinity-ocean/ikakbolit/pkg/errcodes"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func MapSQLError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %v", errcodes.ErrNotFound, err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23502", "23503", "23505", "23514", "22001", "22002", "22003", "22P02", "22007", "42804":
			return fmt.Errorf("%w: %s", errcodes.ErrBadRequest, pgErr.Message)
		}
	}

	return err
}

package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MakePool() (*pgxpool.Pool, error) {
	PG_HOST := os.Getenv("POSTGRES_HOST")
	PG_DB := os.Getenv("POSTGRES_DB")
	PG_USER := os.Getenv("POSTGRES_USER")
	PG_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	PG_PORT := os.Getenv("POSTGRES_PORT")
	PG_SSL := os.Getenv("POSTGRES_SSL")

	missingFields := []string{}

	if PG_HOST == "" {
		missingFields = append(missingFields, "POSTGRES_HOST")
	}
	if PG_DB == "" {
		missingFields = append(missingFields, "POSTGRES_DB")
	}
	if PG_USER == "" {
		missingFields = append(missingFields, "POSTGRES_USER")
	}
	if PG_PASSWORD == "" {
		missingFields = append(missingFields, "POSTGRES_PASSWORD")
	}
	if PG_PORT == "" {
		missingFields = append(missingFields, "POSTGRES_PORT")
	}
	if PG_SSL == "" {
		missingFields = append(missingFields, "POSTGRES_SSL")
	}

	if len(missingFields) > 0 {
		return nil, errors.New("missing fields in .env for db")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		PG_USER,
		PG_PASSWORD,
		PG_HOST,
		PG_PORT,
		PG_DB,
		PG_SSL,
	)

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to create a connection pool: %w", err)
	}
	return dbpool, nil
}
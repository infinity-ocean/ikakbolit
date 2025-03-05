package repo

import (
	"context"
	"fmt"

	"github.com/infinity-ocean/ikakbolit/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MakePool(config config.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.PG_USER,
		config.PG_PASSWORD,
		config.PG_HOST,
		config.PG_PORT,
		config.PG_DB,
		config.PG_SSL,
	)
	
	// Create the connection pool
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to create a connection pool: %w", err)
	}
	return dbpool, nil
}
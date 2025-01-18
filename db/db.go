package db

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

/**
@TODO move from log to slog
*/

/*
  postgres files
*/

//go:embed functions/listings_delta.sql
var listings_delta string

//go:embed functions/tables.sql
var tables string

func NewDatabase(connectionString string) (*pgxpool.Pool, error) {
	//connect to db
	dbpool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		slog.Error("failed to get new pool, likely connection string issue. ")
		return nil, err
	}

	// ping test
	err = dbpool.Ping(context.Background())
	if err != nil {
		slog.Error("Ping test failed.")
		return nil, err
	}

	// create tables
	_, err = dbpool.Exec(context.Background(), tables)
	if err != nil {
		slog.Error("failed to create tables.", slog.Any("error", err))
		return nil, err
	}

	// listings delta function
	_, err = dbpool.Exec(context.Background(), listings_delta)
	if err != nil {
		slog.Error("failed to create listings delta function.", slog.Any("error", err))
		return nil, err
	}

	slog.Info("pool acquired.")

	return dbpool, nil

}

package db

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/jackc/pgx/v5"
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

func NewDatabase(connectionString string) (*pgx.Conn, error) {
	//connect to db
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		slog.Error("failed to get new pool, likely connection string issue. ", slog.Any("error", err))
		return nil, err
	}

	// ping test
	err = conn.Ping(context.Background())
	if err != nil {
		slog.Error("Ping test failed.")
		return nil, err
	}

	// create tables
	_, err = conn.Exec(context.Background(), tables)
	if err != nil {
		slog.Error("failed to create tables.", slog.Any("error", err))
		return nil, err
	}

	// listings delta function
	_, err = conn.Exec(context.Background(), listings_delta)
	if err != nil {
		slog.Error("failed to create listings delta function.", slog.Any("error", err))
		return nil, err
	}

	return conn, nil

}

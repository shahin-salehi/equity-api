package main

import (
	"log/slog"
	"os"

	"github.com/shahin-salehi/equity-api/cmd/api"
	"github.com/shahin-salehi/equity-api/config"
	"github.com/shahin-salehi/equity-api/db"
)

func main() {
	// env

	config := config.InitConfig()
	cs := config.ConnectionString
	if cs == "" {
		slog.Error("connection string empty, shutting down.")
		os.Exit(1)
	}

	db, err := db.NewDatabase(cs)
	if err != nil {
		slog.Error("failed to start database, shutting down.")
		os.Exit(1)
	}

	server := api.NewAPIServer(":5500", db)
	err = server.Run()
	if err != nil {
		slog.Error("Failed to start api server, shutting down.", slog.Any("error", err))
		os.Exit(1)
	}
}

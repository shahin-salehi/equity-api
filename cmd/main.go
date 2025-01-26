package main

import (
	"context"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/shahin-salehi/equity-api/cmd/api"
	"github.com/shahin-salehi/equity-api/config"
	"github.com/shahin-salehi/equity-api/db"
)

func main() {
	// env
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	config := config.InitConfig()
	cs := config.ConnectionString
	if cs == "" {
		slog.Error("connection string empty, shutting down.")
		os.Exit(1)
	}
	slog.Info("connection string log", slog.Any("cs",cs))
	slog.Info("version 1.0.0") 
	
	db, err := db.NewDatabase(cs)
	if err != nil {
		slog.Error("failed to start database, shutting down.")
		os.Exit(1)
	}

	slog.Info("starting server.")

	defer db.Close(context.Background())

	server := api.NewAPIServer(":5500", db)
	err = server.Run()
	if err != nil {
		slog.Error("Failed to start api server, shutting down.", slog.Any("error", err))
		os.Exit(1)
	}
}

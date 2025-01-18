package main

import (
	"log/slog"
	"os"

	"github.com/shahin-salehi/equity-api/cmd/api"
)

func main() {
	server := api.NewAPIServer(":5500", nil)
	err := server.Run()
	if err != nil {
		slog.Error("Failed to start api server, shutting down.", slog.Any("error", err))
		os.Exit(1)
	}
}

package main

import (
	"log/slog"

	"github.com/shahin-salehi/equity-api/api"
)

func main() {
	server := api.NewAPIServer("5500", nil)
	err := server.Run()
	if err != nil {
		slog.Error("Failed to start api server", slog.Any("error", err))
	}
}

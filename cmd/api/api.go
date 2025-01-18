package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/shahin-salehi/equity-api/services/listing"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(port string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: port,
		db:   db,
	}
}

/*
* Starts our api server
 */
func (s *APIServer) Run() error {

	// create router (name mux comes from multiplexer)
	router := http.NewServeMux()

	// give our handler the router so it can register functionality on it
	listingHandler := listing.NewHandler()

	// register desired listing routes on router
	listingHandler.RegisterRoutes(router)

	slog.Info("API server started.")
	return http.ListenAndServe(s.addr, router)
}

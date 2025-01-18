package api

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shahin-salehi/equity-api/services/listing"
)

type APIServer struct {
	addr string
	db   *pgxpool.Pool
}

func NewAPIServer(port string, db *pgxpool.Pool) *APIServer {
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
	listingStore := listing.NewStore(s.db)
	listingHandler := listing.NewHandler(listingStore)

	// register desired listing routes on router
	listingHandler.RegisterRoutes(router)

	slog.Info("API server started.", slog.Any("port", s.addr))
	return http.ListenAndServe(s.addr, router)
}

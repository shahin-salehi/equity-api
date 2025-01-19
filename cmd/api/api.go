package api

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/shahin-salehi/equity-api/db"
	"github.com/shahin-salehi/equity-api/services/listing"
)

type APIServer struct {
	addr string
	db   *pgx.Conn
}

func NewAPIServer(port string, db *pgx.Conn) *APIServer {
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

	// if crud doesn't implement the interface it will get upset here
	dbRepo := db.NewRepo(s.db)
	listingHandler := listing.NewHandler(dbRepo)

	// register desired listing routes on router
	listingHandler.RegisterRoutes(router)
	return http.ListenAndServe(s.addr, router)

}

package listing

import (
	"log/slog"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/listing", h.InsertListing)
}

// post
func (h *Handler) InsertListing(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		slog.Info("Method not allowed", slog.Any("status", 405))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

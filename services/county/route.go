package county

import (
	"log/slog"
	"net/http"

	"github.com/shahin-salehi/equity-api/db"
	"github.com/shahin-salehi/equity-api/utils"
)

// handler wants access to db
type Handler struct {
	db db.CRUD
}

func NewHandler(repo db.CRUD) *Handler {
	return &Handler{db: repo}
}

// add route to our mux
func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/counties", h.Counties)
}

func (h *Handler) Counties(w http.ResponseWriter, r *http.Request) {
	// Reject get
	if r.Method != "GET" {
		slog.Info("Method not allowed", slog.Any("status", 405))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//db call
	counties, err := h.db.ReadCounties()
	if err != nil {
		slog.Error("db retruned error", slog.Any("function", "counties"))
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// ok
	err = utils.SerializeJSON(w, http.StatusOK, counties)
	if err != nil {
		slog.Error("failed to write response", slog.Any("function", "counties"), slog.Any("error", err))
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}

}

package tiles

import (
	"log/slog"
	"net/http"

	"github.com/shahin-salehi/equity-api/db"
	"github.com/shahin-salehi/equity-api/types"
	"github.com/shahin-salehi/equity-api/utils"
)

type Handler struct {
	db db.CRUD
}

// take interface
func NewHandler(repo db.CRUD) *Handler {
	return &Handler{db: repo}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/tiles", h.Readtiles)
	router.HandleFunc("/tiles/insert", h.InsertTiles)
}

func (h *Handler) InsertTiles(w http.ResponseWriter, r *http.Request) {
	// Reject wrong method
	if r.Method != "POST" {
		slog.Info("Method not allowed", slog.Any("status", 405))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// read body
	defer r.Body.Close()
	var b []types.GeoTile
	err := utils.DeserializeJSON(r, &b)
	if err != nil {
		slog.Error("failed to deserialize", slog.Any("error", err))
		utils.WriteError(w, http.StatusUnprocessableEntity, nil)
		return
	}

	// insert
	err = h.db.InsertManyTiles(b)
	if err != nil {
		slog.Error("db returned error when called from handler", slog.Any("function", "InsertTiles"))
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}

	//ok
	err = utils.SerializeJSON(w, http.StatusCreated, nil)
	if err != nil {
		slog.Error("Failed to write response, this shouldn't happen.", slog.Any("function", "InsertTiles"))
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}

}

func (h *Handler) Readtiles(w http.ResponseWriter, r *http.Request) {

	// Reject wrong method
	if r.Method != "GET" {
		slog.Info("Method not allowed", slog.Any("status", 405))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// get param (this panics bad url which is crazy )
	county := r.URL.Query().Get("county")
	if county == "" {
		utils.WriteError(w, http.StatusBadRequest, nil)
	}

	// crud
	tiles, err := h.db.TilesByCounty(county)
	if err != nil {
		slog.Error("failed to get tiles", slog.Any("error", err))
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//ok
	err = utils.SerializeJSON(w, http.StatusOK, tiles)
	if err != nil {
		slog.Error("failed to write tiles", slog.Any("error", err))
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}

}

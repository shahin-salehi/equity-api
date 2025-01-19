package listing

import (
	"io"
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
	router.HandleFunc("/listing", h.InsertListing)
	router.HandleFunc("/delta", h.Delta)
}

func (h *Handler) Delta(w http.ResponseWriter, r *http.Request) {

	// close
	defer r.Body.Close()

	// Reject get
	if r.Method != "POST" {
		slog.Info("Method not allowed", slog.Any("status", 405))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var ids types.ScrapedIDs
	err := utils.DeserializeJSON(r, &ids)
	if err != nil {
		slog.Error("unprocessable entity", slog.Any("status", http.StatusUnprocessableEntity), slog.Any("function", "delta"), slog.Any("error", err))

		// Debug
		b, _ := io.ReadAll(r.Body)
		slog.Debug("body", slog.Any("body", string(b)))

		utils.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	// db
	delta, err := h.db.Delta(ids.IDs)
	if err != nil {
		slog.Error("db retruned error", slog.Any("function", "delta"))
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// ok
	err = utils.SerializeJSON(w, http.StatusOK, delta)
	if err != nil {
		slog.Error("faild to write response, this shouldn't happen", slog.Any("error", err), slog.Any("function", "Delta"))
		return
	}

}

// insert hemnet listing
// TODO: fix retun status for conflict
func (h *Handler) InsertListing(w http.ResponseWriter, r *http.Request) {
	// close
	defer r.Body.Close()

	// Reject get
	if r.Method != "POST" {
		slog.Info("Method not allowed", slog.Any("status", 405))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// get payload
	var payload types.Listing
	err := utils.DeserializeJSON(r, &payload)
	if err != nil {
		slog.Error("unprocessable entity", slog.Any("status", http.StatusUnprocessableEntity))

		// Debug
		b, _ := io.ReadAll(r.Body)
		slog.Debug("body", slog.Any("body", string(b)))

		utils.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	// insert
	err = h.db.Listing(payload)
	if err != nil {
		slog.Error("db retruned error", slog.Any("function", "InsertListing"))
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// ok
	slog.Info("created", slog.Any("status", http.StatusCreated))
	err = utils.SerializeJSON(w, http.StatusCreated, nil)
	if err != nil {
		slog.Error("faild to write response, this shouldn't happen", slog.Any("error", err), slog.Any("function", "InsertListing"))
		return
	}

}

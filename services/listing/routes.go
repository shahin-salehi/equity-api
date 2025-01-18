package listing

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/shahin-salehi/equity-api/types"
	"github.com/shahin-salehi/equity-api/utils"
)

type Handler struct {
	store types.ListingStore
}

// take interface
func NewHandler(store types.ListingStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/listing", h.InsertListing)
}

// insert hemnet listing
// TODO: fix retun status for conflict
func (h *Handler) InsertListing(w http.ResponseWriter, r *http.Request) {
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
		defer r.Body.Close()
		b, _ := io.ReadAll(r.Body)
		slog.Debug("body", slog.Any("body", string(b)))

		utils.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	// insert
	err = h.store.Listing(payload)
	if err != nil {
		slog.Error("failed to insert listing", slog.Any("error", err))
		http.Error(w, "database returned error", http.StatusInternalServerError)
		return
	}

	// ok
	slog.Info("created", slog.Any("status", http.StatusCreated))
	utils.SerializeJSON(w, http.StatusCreated, nil)

}

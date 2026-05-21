package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/barberku/backend-barber/internal/service"
	"github.com/barberku/backend-barber/pkg/response"
)

type HistoryHandlerImpl struct {
	historyService service.HistoryService
}

func NewHistoryHandler(historyService service.HistoryService) *HistoryHandlerImpl {
	return &HistoryHandlerImpl{
		historyService: historyService,
	}
}

func (h *HistoryHandlerImpl) RegisterRoutes(r chi.Router) {
	r.Get("/", h.GetHistory)
	r.Post("/{id}/rate", h.RateService)
}

func (h *HistoryHandlerImpl) GetHistory(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		response.BadRequest(w, "Customer ID wajib diisi")
		return
	}

	history, err := h.historyService.GetHistory(r.Context(), customerID)
	if err != nil {
		slog.Error("failed to get history", "error", err)
		response.InternalServerError(w, "Gagal mengambil riwayat")
		return
	}

	response.Success(w, "OK", history)
}

func (h *HistoryHandlerImpl) RateService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	if err := h.historyService.RateService(r.Context(), id, req.Rating, req.Comment); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Terima kasih atas rating Anda", nil)
}

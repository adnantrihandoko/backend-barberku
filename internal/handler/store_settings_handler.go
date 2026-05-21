package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/service"
	"github.com/barberku/backend-barber/pkg/response"
)

type StoreSettingsHandlerImpl struct {
	settingsService service.StoreSettingsService
}

func NewStoreSettingsHandler(settingsService service.StoreSettingsService) *StoreSettingsHandlerImpl {
	return &StoreSettingsHandlerImpl{
		settingsService: settingsService,
	}
}

func (h *StoreSettingsHandlerImpl) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /api/v1/settings", h.GetSettings)
	r.HandleFunc("PUT /api/v1/settings", h.UpdateSettings)
}

func (h *StoreSettingsHandlerImpl) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settingsService.GetSettings(r.Context())
	if err != nil {
		slog.Error("failed to get settings", "error", err)
		response.InternalServerError(w, "Gagal mengambil pengaturan")
		return
	}

	response.Success(w, "OK", settings)
}

func (h *StoreSettingsHandlerImpl) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OpenHour     int `json:"open_hour"`
		CloseHour    int `json:"close_hour"`
		MaxQueueSize int `json:"max_queue_size"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	settings := &entity.StoreSettings{
		OpenHour:     req.OpenHour,
		CloseHour:    req.CloseHour,
		MaxQueueSize: req.MaxQueueSize,
	}

	if err := h.settingsService.UpdateSettings(r.Context(), settings); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, "Pengaturan berhasil diperbarui", settings)
}

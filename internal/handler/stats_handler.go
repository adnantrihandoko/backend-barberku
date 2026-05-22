package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/barberku/backend-barber/internal/service"
	"github.com/barberku/backend-barber/pkg/response"
)

type StatsHandlerImpl struct {
	statsService service.StatsService
}

func NewStatsHandler(statsService service.StatsService) *StatsHandlerImpl {
	return &StatsHandlerImpl{
		statsService: statsService,
	}
}

func (h *StatsHandlerImpl) RegisterRoutes(r chi.Router) {
	r.Get("/", h.GetDailyStats)
}

func (h *StatsHandlerImpl) GetDailyStats(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	
	var date time.Time
	var err error
	
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			response.BadRequest(w, "Format tanggal tidak valid (YYYY-MM-DD)")
			return
		}
	} else {
		date = time.Now()
	}

	stats, err := h.statsService.GetDailyStats(r.Context(), date)
	if err != nil {
		slog.Error("failed to get stats", "error", err)
		response.InternalServerError(w, "Gagal mengambil statistik")
		return
	}

	response.Success(w, "OK", stats)
}

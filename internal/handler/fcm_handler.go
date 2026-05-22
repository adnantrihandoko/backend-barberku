package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/barberku/backend-barber/internal/service"
	"github.com/barberku/backend-barber/pkg/response"
)

type FCMHandlerImpl struct {
	fcmService service.FCMService
}

func NewFCMHandler(fcmService service.FCMService) *FCMHandlerImpl {
	return &FCMHandlerImpl{
		fcmService: fcmService,
	}
}

func (h *FCMHandlerImpl) RegisterRoutes(r chi.Router) {
	r.Post("/register", h.RegisterToken)
}

type registerTokenRequest struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

func (h *FCMHandlerImpl) RegisterToken(w http.ResponseWriter, r *http.Request) {
	var req registerTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Request tidak valid")
		return
	}

	if req.Token == "" || req.Platform == "" {
		response.BadRequest(w, "Token dan platform wajib diisi")
		return
	}

	customerID := r.Context().Value("customer_id")
	if customerID == nil {
		response.Unauthorized(w, "Customer ID tidak ditemukan")
		return
	}

	if err := h.fcmService.RegisterToken(r.Context(), customerID.(string), req.Token, req.Platform); err != nil {
		slog.Error("failed to register FCM token", "error", err)
		response.InternalServerError(w, "Gagal mendaftarkan token FCM")
		return
	}

	response.Success(w, "Token FCM berhasil didaftarkan", nil)
}

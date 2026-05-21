package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/service"
	"github.com/barberku/backend-barber/pkg/response"
)

type AuthHandlerImpl struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		authService: authService,
	}
}

func (h *AuthHandlerImpl) RegisterRoutes(r chi.Router) {
	r.Post("/login", h.Login)
	r.Get("/me", h.GetMe)
}

func (h *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	var req entity.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	loginResp, err := h.authService.Login(r.Context(), req.Email, req.PIN)
	if err != nil {
		slog.Warn("login failed", "email", req.Email, "error", err)
		response.Unauthorized(w, err.Error())
		return
	}

	slog.Info("login successful", "email", req.Email)
	response.Success(w, "Login berhasil", loginResp)
}

func (h *AuthHandlerImpl) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*entity.JWTClaims)
	if !ok {
		response.Unauthorized(w, "Token tidak valid")
		return
	}

	response.Success(w, "OK", claims)
}

package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/barberku/backend-barber/internal/service"
	"github.com/barberku/backend-barber/pkg/response"
)

type BarberHandlerImpl struct {
	barberService service.BarberService
}

func NewBarberHandler(barberService service.BarberService) *BarberHandlerImpl {
	return &BarberHandlerImpl{
		barberService: barberService,
	}
}

func (h *BarberHandlerImpl) RegisterRoutes(r chi.Router) {
	r.Get("/", h.ListBarbers)
	r.Get("/{id}", h.GetBarber)
	r.Post("/", h.CreateBarber)
	r.Put("/{id}", h.UpdateBarber)
	r.Delete("/{id}", h.DeleteBarber)
}

func (h *BarberHandlerImpl) ListBarbers(w http.ResponseWriter, r *http.Request) {
	barbers, err := h.barberService.ListBarbers(r.Context())
	if err != nil {
		slog.Error("failed to list barbers", "error", err)
		response.InternalServerError(w, "Gagal mengambil data barber")
		return
	}

	response.Success(w, "OK", barbers)
}

func (h *BarberHandlerImpl) GetBarber(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	barber, err := h.barberService.GetBarber(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Barber tidak ditemukan")
		return
	}

	response.Success(w, "OK", barber)
}

func (h *BarberHandlerImpl) CreateBarber(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      string `json:"name"`
		Specialty string `json:"specialty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "Nama barber wajib diisi")
		return
	}

	newBarber, err := h.barberService.CreateBarber(r.Context(), req.Name, req.Specialty)
	if err != nil {
		slog.Error("failed to create barber", "error", err)
		response.InternalServerError(w, "Gagal menambah barber")
		return
	}

	response.Created(w, "Barber berhasil ditambahkan", newBarber)
}

func (h *BarberHandlerImpl) UpdateBarber(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Name      string `json:"name"`
		Specialty string `json:"specialty"`
		IsActive  bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "Nama barber wajib diisi")
		return
	}

	updatedBarber, err := h.barberService.UpdateBarber(r.Context(), id, req.Name, req.Specialty, req.IsActive)
	if err != nil {
		handleBarberError(w, err)
		return
	}

	response.Success(w, "Barber berhasil diperbarui", updatedBarber)
}

func (h *BarberHandlerImpl) DeleteBarber(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.barberService.DeleteBarber(r.Context(), id); err != nil {
		handleBarberError(w, err)
		return
	}

	response.Success(w, "Barber berhasil dihapus", nil)
}

func handleBarberError(w http.ResponseWriter, err error) {
	msg := err.Error()
	if strings.Contains(msg, "tidak ditemukan") {
		response.NotFound(w, msg)
	} else {
		slog.Error("barber operation failed", "error", err)
		response.InternalServerError(w, "Operasi barber gagal")
	}
}

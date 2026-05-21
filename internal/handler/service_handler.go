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

type ServiceHandlerImpl struct {
	serviceService service.ServiceService
}

func NewServiceHandler(serviceService service.ServiceService) *ServiceHandlerImpl {
	return &ServiceHandlerImpl{
		serviceService: serviceService,
	}
}

func (h *ServiceHandlerImpl) RegisterRoutes(r chi.Router) {
	r.Get("/", h.ListServices)
	r.Get("/{id}", h.GetService)
	r.Post("/", h.CreateService)
	r.Put("/{id}", h.UpdateService)
	r.Delete("/{id}", h.DeleteService)
}

func (h *ServiceHandlerImpl) ListServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.serviceService.ListServices(r.Context())
	if err != nil {
		slog.Error("failed to list services", "error", err)
		response.InternalServerError(w, "Gagal mengambil data layanan")
		return
	}

	response.Success(w, "OK", services)
}

func (h *ServiceHandlerImpl) GetService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	service, err := h.serviceService.GetService(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Layanan tidak ditemukan")
		return
	}

	response.Success(w, "OK", service)
}

func (h *ServiceHandlerImpl) CreateService(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Duration    int     `json:"duration"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "Nama layanan wajib diisi")
		return
	}

	newService, err := h.serviceService.CreateService(r.Context(), req.Name, req.Description, req.Price, req.Duration)
	if err != nil {
		slog.Error("failed to create service", "error", err)
		response.InternalServerError(w, "Gagal menambah layanan")
		return
	}

	response.Created(w, "Layanan berhasil ditambahkan", newService)
}

func (h *ServiceHandlerImpl) UpdateService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Duration    int     `json:"duration"`
		IsActive    bool    `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	if req.Name == "" {
		response.BadRequest(w, "Nama layanan wajib diisi")
		return
	}

	updatedService, err := h.serviceService.UpdateService(r.Context(), id, req.Name, req.Description, req.Price, req.Duration, req.IsActive)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, "Layanan berhasil diperbarui", updatedService)
}

func (h *ServiceHandlerImpl) DeleteService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.serviceService.DeleteService(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	response.Success(w, "Layanan berhasil dihapus", nil)
}

func handleServiceError(w http.ResponseWriter, err error) {
	msg := err.Error()
	if strings.Contains(msg, "tidak ditemukan") {
		response.NotFound(w, msg)
	} else {
		slog.Error("service operation failed", "error", err)
		response.InternalServerError(w, "Operasi layanan gagal")
	}
}

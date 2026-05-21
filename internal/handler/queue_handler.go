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

type QueueHandlerImpl struct {
	queueService service.QueueService
}

func NewQueueHandler(queueService service.QueueService) *QueueHandlerImpl {
	return &QueueHandlerImpl{
		queueService: queueService,
	}
}

func (h *QueueHandlerImpl) RegisterRoutes(r chi.Router) {
	r.Get("/", h.GetQueueList)
	r.Post("/join", h.JoinQueue)
	r.Post("/walk-in", h.AddWalkIn)
	r.Get("/{id}", h.GetQueueDetail)
	r.Post("/{id}/call", h.CallQueue)
	r.Post("/{id}/serve", h.ServeQueue)
	r.Post("/{id}/complete", h.CompleteQueue)
	r.Post("/{id}/skip", h.SkipQueue)
	r.Post("/{id}/cancel", h.CancelQueue)
}

func (h *QueueHandlerImpl) GetQueueList(w http.ResponseWriter, r *http.Request) {
	queues, err := h.queueService.GetQueueList(r.Context())
	if err != nil {
		slog.Error("failed to get queue list", "error", err)
		response.InternalServerError(w, "Gagal mengambil data antrian")
		return
	}

	response.Success(w, "OK", queues)
}

func (h *QueueHandlerImpl) GetQueueDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	queue, err := h.queueService.GetQueueDetail(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Antrian tidak ditemukan")
		return
	}

	response.Success(w, "OK", queue)
}

func (h *QueueHandlerImpl) JoinQueue(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID   string  `json:"customer_id"`
		CustomerName string  `json:"customer_name"`
		ServiceID    string  `json:"service_id"`
		ServiceName  string  `json:"service_name"`
		BarberID     *string `json:"barber_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	if req.CustomerID == "" || req.ServiceID == "" {
		response.BadRequest(w, "Customer ID dan Service ID wajib diisi")
		return
	}

	queue, err := h.queueService.JoinQueue(r.Context(), req.CustomerID, req.CustomerName, req.ServiceID, req.ServiceName, req.BarberID)
	if err != nil {
		handleJoinError(w, err)
		return
	}

	response.Created(w, "Berhasil bergabung ke antrian", queue)
}

func (h *QueueHandlerImpl) AddWalkIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerName string  `json:"customer_name"`
		ServiceID    string  `json:"service_id"`
		ServiceName  string  `json:"service_name"`
		BarberID     *string `json:"barber_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	if req.CustomerName == "" || req.ServiceID == "" {
		response.BadRequest(w, "Nama pelanggan dan Service ID wajib diisi")
		return
	}

	queue, err := h.queueService.AddWalkIn(r.Context(), req.CustomerName, req.ServiceID, req.ServiceName, req.BarberID)
	if err != nil {
		slog.Error("failed to add walk-in", "error", err)
		response.InternalServerError(w, "Gagal menambah antrian walk-in")
		return
	}

	response.Created(w, "Berhasil menambah antrian walk-in", queue)
}

func (h *QueueHandlerImpl) CallQueue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.queueService.CallQueue(r.Context(), id); err != nil {
		handleQueueError(w, err)
		return
	}

	response.Success(w, "Antrian berhasil dipanggil", nil)
}

func (h *QueueHandlerImpl) ServeQueue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.queueService.ServeQueue(r.Context(), id); err != nil {
		handleQueueError(w, err)
		return
	}

	response.Success(w, "Layanan dimulai", nil)
}

func (h *QueueHandlerImpl) CompleteQueue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.queueService.CompleteQueue(r.Context(), id); err != nil {
		handleQueueError(w, err)
		return
	}

	response.Success(w, "Layanan selesai", nil)
}

func (h *QueueHandlerImpl) SkipQueue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.queueService.SkipQueue(r.Context(), id); err != nil {
		handleQueueError(w, err)
		return
	}

	response.Success(w, "Antrian dilewati (no-show)", nil)
}

func (h *QueueHandlerImpl) CancelQueue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.queueService.CancelQueue(r.Context(), id); err != nil {
		handleQueueError(w, err)
		return
	}

	response.Success(w, "Antrian dibatalkan", nil)
}

func handleQueueError(w http.ResponseWriter, err error) {
	msg := err.Error()
	if strings.Contains(msg, "tidak ditemukan") {
		response.NotFound(w, msg)
	} else if strings.Contains(msg, "tidak dalam status") || strings.Contains(msg, "sudah selesai") {
		response.BadRequest(w, msg)
	} else {
		slog.Error("queue operation failed", "error", err)
		response.InternalServerError(w, "Operasi antrian gagal")
	}
}

func handleJoinError(w http.ResponseWriter, err error) {
	msg := err.Error()
	if strings.Contains(msg, "jam operasional") {
		response.BadRequest(w, "Maaf, barbershop sedang tutup. Jam operasional: 09:00 - 21:00")
	} else if strings.Contains(msg, "antrian sudah penuh") {
		response.BadRequest(w, "Antrian hari ini sudah penuh. Silakan datang besok.")
	} else if strings.Contains(msg, "antrian aktif") {
		response.BadRequest(w, "Anda sudah memiliki antrian aktif. Selesaikan antrian terlebih dahulu.")
	} else if strings.Contains(msg, "pembatalan") {
		response.Forbidden(w, msg)
	} else {
		slog.Error("failed to join queue", "error", err)
		response.InternalServerError(w, "Gagal bergabung ke antrian")
	}
}

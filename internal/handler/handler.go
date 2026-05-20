package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type QueueHandler interface {
	RegisterRoutes(r chi.Router)
	JoinQueue(w http.ResponseWriter, r *http.Request)
	CancelQueue(w http.ResponseWriter, r *http.Request)
	GetQueueList(w http.ResponseWriter, r *http.Request)
	GetQueueDetail(w http.ResponseWriter, r *http.Request)
	CallQueue(w http.ResponseWriter, r *http.Request)
	CompleteQueue(w http.ResponseWriter, r *http.Request)
	SkipQueue(w http.ResponseWriter, r *http.Request)
	AddWalkIn(w http.ResponseWriter, r *http.Request)
}

type AuthHandler interface {
	RegisterRoutes(r chi.Router)
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}

type BarberHandler interface {
	RegisterRoutes(r chi.Router)
	ListBarbers(w http.ResponseWriter, r *http.Request)
	GetBarber(w http.ResponseWriter, r *http.Request)
	CreateBarber(w http.ResponseWriter, r *http.Request)
	UpdateBarber(w http.ResponseWriter, r *http.Request)
	DeleteBarber(w http.ResponseWriter, r *http.Request)
}

type ServiceHandler interface {
	RegisterRoutes(r chi.Router)
	ListServices(w http.ResponseWriter, r *http.Request)
	GetService(w http.ResponseWriter, r *http.Request)
	CreateService(w http.ResponseWriter, r *http.Request)
	UpdateService(w http.ResponseWriter, r *http.Request)
	DeleteService(w http.ResponseWriter, r *http.Request)
}

type WebSocketHandler interface {
	RegisterRoutes(r chi.Router)
	HandleWebSocket(w http.ResponseWriter, r *http.Request)
}

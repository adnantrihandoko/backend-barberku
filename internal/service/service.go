package service

import (
	"context"

	"github.com/barberku/backend-barber/internal/entity"
)

type QueueService interface {
	JoinQueue(ctx context.Context, customerID, serviceID string, barberID *string) (*entity.Queue, error)
	CancelQueue(ctx context.Context, queueID string) error
	CallQueue(ctx context.Context, queueID string) error
	CompleteQueue(ctx context.Context, queueID string) error
	SkipQueue(ctx context.Context, queueID string) error
	GetQueueList(ctx context.Context) ([]entity.Queue, error)
	GetQueueDetail(ctx context.Context, queueID string) (*entity.Queue, error)
	AddWalkIn(ctx context.Context, customerName, serviceID string, barberID *string) (*entity.Queue, error)
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, *entity.User, error)
	Register(ctx context.Context, name, email, phone, password string, role entity.UserRole) (*entity.User, error)
	ValidateToken(ctx context.Context, token string) (*entity.User, error)
}

type BarberService interface {
	ListBarbers(ctx context.Context) ([]entity.Barber, error)
	GetBarber(ctx context.Context, id string) (*entity.Barber, error)
	CreateBarber(ctx context.Context, name, specialty string) (*entity.Barber, error)
	UpdateBarber(ctx context.Context, id, name, specialty string, isActive bool) (*entity.Barber, error)
	DeleteBarber(ctx context.Context, id string) error
}

type ServiceService interface {
	ListServices(ctx context.Context) ([]entity.Service, error)
	GetService(ctx context.Context, id string) (*entity.Service, error)
	CreateService(ctx context.Context, name, description string, price float64, duration int) (*entity.Service, error)
	UpdateService(ctx context.Context, id, name, description string, price float64, duration int, isActive bool) (*entity.Service, error)
	DeleteService(ctx context.Context, id string) error
}

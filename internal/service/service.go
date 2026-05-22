package service

import (
	"context"
	"time"

	"github.com/barberku/backend-barber/internal/entity"
)

type AuthService interface {
	Login(ctx context.Context, email, pin string) (*entity.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*entity.JWTClaims, error)
}

type QueueService interface {
	GetQueueList(ctx context.Context) ([]entity.Queue, error)
	GetQueueDetail(ctx context.Context, queueID string) (*entity.Queue, error)
	JoinQueue(ctx context.Context, customerID, customerName, serviceID, serviceName string, barberID *string) (*entity.Queue, error)
	CallQueue(ctx context.Context, queueID string) error
	ServeQueue(ctx context.Context, queueID string) error
	CompleteQueue(ctx context.Context, queueID string) error
	SkipQueue(ctx context.Context, queueID string) error
	CancelQueue(ctx context.Context, queueID string) error
	AddWalkIn(ctx context.Context, customerName, serviceID, serviceName string, barberID *string) (*entity.Queue, error)
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

type HistoryService interface {
	GetHistory(ctx context.Context, customerID string) ([]entity.Queue, error)
	RateService(ctx context.Context, queueID string, rating int, comment string) error
}

type StatsService interface {
	GetDailyStats(ctx context.Context, date time.Time) (*entity.DailyStats, error)
}

type StoreSettingsService interface {
	GetSettings(ctx context.Context) (*entity.StoreSettings, error)
	UpdateSettings(ctx context.Context, settings *entity.StoreSettings) error
}

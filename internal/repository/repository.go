package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/barberku/backend-barber/internal/entity"
)

type QueueRepository interface {
	Create(ctx context.Context, queue *entity.Queue) error
	GetByID(ctx context.Context, id string) (*entity.Queue, error)
	GetList(ctx context.Context) ([]entity.Queue, error)
	Update(ctx context.Context, queue *entity.Queue) error
	UpdateStatus(ctx context.Context, id string, status entity.QueueStatus) error
	Delete(ctx context.Context, id string) error
	GetByCustomerID(ctx context.Context, customerID string) ([]entity.Queue, error)
	GetCountByStatus(ctx context.Context, status entity.QueueStatus) (int, error)
	GetNextQueueNumber(ctx context.Context) (int, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	GetCancelCountByCustomerToday(ctx context.Context, customerID string) (int, error)
	GetLastCancelTimeByCustomer(ctx context.Context, customerID string) (*time.Time, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]entity.User, error)
}

type BarberRepository interface {
	Create(ctx context.Context, barber *entity.Barber) error
	GetByID(ctx context.Context, id string) (*entity.Barber, error)
	List(ctx context.Context) ([]entity.Barber, error)
	Update(ctx context.Context, barber *entity.Barber) error
	Delete(ctx context.Context, id string) error
}

type ServiceRepository interface {
	Create(ctx context.Context, service *entity.Service) error
	GetByID(ctx context.Context, id string) (*entity.Service, error)
	List(ctx context.Context) ([]entity.Service, error)
	Update(ctx context.Context, service *entity.Service) error
	Delete(ctx context.Context, id string) error
}

type FCMTokenRepository interface {
	Save(ctx context.Context, token *entity.FCMToken) error
	GetByCustomerID(ctx context.Context, customerID string) ([]entity.FCMToken, error)
	Delete(ctx context.Context, id string) error
}

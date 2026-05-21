package service

import (
	"context"

	"github.com/barberku/backend-barber/internal/entity"
)

type AuthService interface {
	Login(ctx context.Context, email, pin string) (*entity.LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*entity.JWTClaims, error)
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

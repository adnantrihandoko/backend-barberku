package repository

import (
	"context"

	"github.com/barberku/backend-barber/internal/entity"
)

type AuthRepository interface {
	GetAdminByEmail(ctx context.Context, email string) (*entity.AdminUser, error)
	GetAdminByID(ctx context.Context, id string) (*entity.AdminUser, error)
	CreateAdmin(ctx context.Context, admin *entity.AdminUser) error
	UpdateAdmin(ctx context.Context, admin *entity.AdminUser) error
}

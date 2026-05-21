package repository

import (
	"context"

	"github.com/barberku/backend-barber/internal/entity"
)

type StoreSettingsRepository interface {
	Get(ctx context.Context) (*entity.StoreSettings, error)
	Update(ctx context.Context, settings *entity.StoreSettings) error
}

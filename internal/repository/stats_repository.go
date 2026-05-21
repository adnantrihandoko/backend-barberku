package repository

import (
	"context"
	"time"

	"github.com/barberku/backend-barber/internal/entity"
)

type StatsRepository interface {
	GetDailyStats(ctx context.Context, date time.Time) (*entity.DailyStats, error)
}

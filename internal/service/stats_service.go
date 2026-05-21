package service

import (
	"context"
	"time"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/repository"
)

type StatsServiceImpl struct {
	statsRepo repository.StatsRepository
}

func NewStatsService(statsRepo repository.StatsRepository) *StatsServiceImpl {
	return &StatsServiceImpl{
		statsRepo: statsRepo,
	}
}

func (s *StatsServiceImpl) GetDailyStats(ctx context.Context, date time.Time) (*entity.DailyStats, error) {
	return s.statsRepo.GetDailyStats(ctx, date)
}

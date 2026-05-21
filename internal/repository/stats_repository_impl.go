package repository

import (
	"context"
	"time"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewStatsRepository(pool *pgxpool.Pool) *StatsRepositoryImpl {
	return &StatsRepositoryImpl{Pool: pool}
}

func (r *StatsRepositoryImpl) GetDailyStats(ctx context.Context, date time.Time) (*entity.DailyStats, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE status = 'completed') as total_served,
			COUNT(*) FILTER (WHERE status = 'canceled') as total_canceled,
			COALESCE(AVG(EXTRACT(EPOCH FROM (called_at - created_at)) / 60) FILTER (WHERE called_at IS NOT NULL), 0) as avg_wait_time_min,
			COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - called_at)) / 60) FILTER (WHERE completed_at IS NOT NULL AND called_at IS NOT NULL), 0) as avg_service_time_min,
			COALESCE(SUM(s.price) FILTER (WHERE q.status = 'completed'), 0) as total_revenue
		FROM queues q
		LEFT JOIN services s ON q.service_id = s.id
		WHERE DATE(q.created_at) = $1
	`

	var stats entity.DailyStats
	err := r.Pool.QueryRow(ctx, query, date).Scan(
		&stats.TotalServed,
		&stats.TotalCanceled,
		&stats.AvgWaitTimeMin,
		&stats.AvgServiceTimeMin,
		&stats.TotalRevenue,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

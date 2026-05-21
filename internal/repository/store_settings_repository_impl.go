package repository

import (
	"context"
	"fmt"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreSettingsRepositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewStoreSettingsRepository(pool *pgxpool.Pool) *StoreSettingsRepositoryImpl {
	return &StoreSettingsRepositoryImpl{Pool: pool}
}

func (r *StoreSettingsRepositoryImpl) Get(ctx context.Context) (*entity.StoreSettings, error) {
	query := `
		SELECT open_hour, close_hour, max_queue_size
		FROM store_settings
		ORDER BY id DESC LIMIT 1
	`

	var settings entity.StoreSettings
	err := r.Pool.QueryRow(ctx, query).Scan(
		&settings.OpenHour,
		&settings.CloseHour,
		&settings.MaxQueueSize,
	)

	if err == pgx.ErrNoRows {
		return &entity.StoreSettings{
			OpenHour:     9,
			CloseHour:    21,
			MaxQueueSize: 50,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (r *StoreSettingsRepositoryImpl) Update(ctx context.Context, settings *entity.StoreSettings) error {
	query := `
		UPDATE store_settings
		SET open_hour = $1, close_hour = $2, max_queue_size = $3, updated_at = NOW()
		WHERE id = (SELECT id FROM store_settings ORDER BY id DESC LIMIT 1)
	`

	result, err := r.Pool.Exec(ctx, query,
		settings.OpenHour,
		settings.CloseHour,
		settings.MaxQueueSize,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		query = `
			INSERT INTO store_settings (open_hour, close_hour, max_queue_size, updated_at)
			VALUES ($1, $2, $3, NOW())
		`
		_, err = r.Pool.Exec(ctx, query,
			settings.OpenHour,
			settings.CloseHour,
			settings.MaxQueueSize,
		)
		return err
	}

	return nil
}

package repository

import (
	"context"
	"fmt"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewServiceRepository(pool *pgxpool.Pool) *ServiceRepositoryImpl {
	return &ServiceRepositoryImpl{Pool: pool}
}

func (r *ServiceRepositoryImpl) Create(ctx context.Context, service *entity.Service) error {
	query := `
		INSERT INTO services (id, name, description, price, duration, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.Pool.Exec(ctx, query,
		service.ID,
		service.Name,
		service.Description,
		service.Price,
		service.Duration,
		service.IsActive,
		service.CreatedAt,
		service.UpdatedAt,
	)

	return err
}

func (r *ServiceRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Service, error) {
	query := `
		SELECT id, name, description, price, duration, is_active, created_at, updated_at
		FROM services WHERE id = $1
	`

	var service entity.Service
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		&service.Description,
		&service.Price,
		&service.Duration,
		&service.IsActive,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("service not found")
	}
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *ServiceRepositoryImpl) List(ctx context.Context) ([]entity.Service, error) {
	query := `
		SELECT id, name, description, price, duration, is_active, created_at, updated_at
		FROM services ORDER BY name ASC
	`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []entity.Service
	for rows.Next() {
		var service entity.Service
		if err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Description,
			&service.Price,
			&service.Duration,
			&service.IsActive,
			&service.CreatedAt,
			&service.UpdatedAt,
		); err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, rows.Err()
}

func (r *ServiceRepositoryImpl) Update(ctx context.Context, service *entity.Service) error {
	query := `
		UPDATE services
		SET name = $1, description = $2, price = $3, duration = $4, is_active = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.Pool.Exec(ctx, query,
		service.Name,
		service.Description,
		service.Price,
		service.Duration,
		service.IsActive,
		service.UpdatedAt,
		service.ID,
	)

	return err
}

func (r *ServiceRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM services WHERE id = $1`
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

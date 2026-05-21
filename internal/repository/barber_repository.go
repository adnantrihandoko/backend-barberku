package repository

import (
	"context"
	"fmt"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BarberRepositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewBarberRepository(pool *pgxpool.Pool) *BarberRepositoryImpl {
	return &BarberRepositoryImpl{Pool: pool}
}

func (r *BarberRepositoryImpl) Create(ctx context.Context, barber *entity.Barber) error {
	query := `
		INSERT INTO barbers (id, name, specialty, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.Pool.Exec(ctx, query,
		barber.ID,
		barber.Name,
		barber.Specialty,
		barber.IsActive,
		barber.CreatedAt,
		barber.UpdatedAt,
	)

	return err
}

func (r *BarberRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Barber, error) {
	query := `
		SELECT id, name, specialty, is_active, created_at, updated_at
		FROM barbers WHERE id = $1
	`

	var barber entity.Barber
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&barber.ID,
		&barber.Name,
		&barber.Specialty,
		&barber.IsActive,
		&barber.CreatedAt,
		&barber.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("barber not found")
	}
	if err != nil {
		return nil, err
	}

	return &barber, nil
}

func (r *BarberRepositoryImpl) List(ctx context.Context) ([]entity.Barber, error) {
	query := `
		SELECT id, name, specialty, is_active, created_at, updated_at
		FROM barbers ORDER BY name ASC
	`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var barbers []entity.Barber
	for rows.Next() {
		var barber entity.Barber
		if err := rows.Scan(
			&barber.ID,
			&barber.Name,
			&barber.Specialty,
			&barber.IsActive,
			&barber.CreatedAt,
			&barber.UpdatedAt,
		); err != nil {
			return nil, err
		}
		barbers = append(barbers, barber)
	}

	return barbers, rows.Err()
}

func (r *BarberRepositoryImpl) Update(ctx context.Context, barber *entity.Barber) error {
	query := `
		UPDATE barbers
		SET name = $1, specialty = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.Pool.Exec(ctx, query,
		barber.Name,
		barber.Specialty,
		barber.IsActive,
		barber.UpdatedAt,
		barber.ID,
	)

	return err
}

func (r *BarberRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM barbers WHERE id = $1`
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

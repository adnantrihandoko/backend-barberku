package repository

import (
	"context"
	"fmt"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{Pool: pool}
}

func (r *AuthRepositoryImpl) GetAdminByEmail(ctx context.Context, email string) (*entity.AdminUser, error) {
	query := `
		SELECT id, name, email, phone, role, pin_hash, is_active, created_at, updated_at
		FROM users WHERE email = $1 AND role IN ('admin', 'barber')
	`

	var admin entity.AdminUser
	err := r.Pool.QueryRow(ctx, query, email).Scan(
		&admin.ID,
		&admin.Name,
		&admin.Email,
		&admin.Phone,
		&admin.Role,
		&admin.PINHash,
		&admin.IsActive,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("admin not found")
	}
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *AuthRepositoryImpl) GetAdminByID(ctx context.Context, id string) (*entity.AdminUser, error) {
	query := `
		SELECT id, name, email, phone, role, pin_hash, is_active, created_at, updated_at
		FROM users WHERE id = $1 AND role IN ('admin', 'barber')
	`

	var admin entity.AdminUser
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&admin.ID,
		&admin.Name,
		&admin.Email,
		&admin.Phone,
		&admin.Role,
		&admin.PINHash,
		&admin.IsActive,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("admin not found")
	}
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func (r *AuthRepositoryImpl) CreateAdmin(ctx context.Context, admin *entity.AdminUser) error {
	query := `
		INSERT INTO users (id, name, email, phone, role, pin_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.Pool.Exec(ctx, query,
		admin.ID,
		admin.Name,
		admin.Email,
		admin.Phone,
		admin.Role,
		admin.PINHash,
		admin.IsActive,
		admin.CreatedAt,
		admin.UpdatedAt,
	)

	return err
}

func (r *AuthRepositoryImpl) UpdateAdmin(ctx context.Context, admin *entity.AdminUser) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, phone = $3, role = $4, pin_hash = $5, is_active = $6, updated_at = $7
		WHERE id = $8
	`

	_, err := r.Pool.Exec(ctx, query,
		admin.Name,
		admin.Email,
		admin.Phone,
		admin.Role,
		admin.PINHash,
		admin.IsActive,
		admin.UpdatedAt,
		admin.ID,
	)

	return err
}

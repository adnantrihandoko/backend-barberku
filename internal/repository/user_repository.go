package repository

import (
	"context"
	"fmt"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepositoryImpl {
	return &UserRepositoryImpl{Pool: pool}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, name, email, phone, role, pin_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.Pool.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Phone,
		user.Role,
		"",
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT id, name, email, phone, role, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user entity.User
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, name, email, phone, role, pin_hash, created_at, updated_at
		FROM users WHERE email = $1
	`

	var user entity.User
	var pinHash string
	err := r.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Role,
		&pinHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	query := `
		SELECT id, name, email, phone, role, created_at, updated_at
		FROM users WHERE phone = $1
	`

	var user entity.User
	err := r.Pool.QueryRow(ctx, query, phone).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, phone = $3, role = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.Pool.Exec(ctx, query,
		user.Name,
		user.Email,
		user.Phone,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

func (r *UserRepositoryImpl) List(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT id, name, email, phone, role, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserRepositoryImpl) GetPINHash(ctx context.Context, email string) (string, error) {
	query := `SELECT pin_hash FROM users WHERE email = $1`
	var pinHash string
	err := r.Pool.QueryRow(ctx, query, email).Scan(&pinHash)
	if err == pgx.ErrNoRows {
		return "", fmt.Errorf("user not found")
	}
	return pinHash, err
}

func (r *UserRepositoryImpl) UpdatePINHash(ctx context.Context, email, pinHash string) error {
	query := `UPDATE users SET pin_hash = $1, updated_at = NOW() WHERE email = $2`
	_, err := r.Pool.Exec(ctx, query, pinHash, email)
	return err
}

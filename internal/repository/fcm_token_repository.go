package repository

import (
	"context"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FCMTokenRepositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewFCMTokenRepository(pool *pgxpool.Pool) *FCMTokenRepositoryImpl {
	return &FCMTokenRepositoryImpl{Pool: pool}
}

func (r *FCMTokenRepositoryImpl) Save(ctx context.Context, token *entity.FCMToken) error {
	query := `
		INSERT INTO fcm_tokens (id, customer_id, token, platform, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (customer_id, token)
		DO UPDATE SET platform = $4, updated_at = $6
	`

	_, err := r.Pool.Exec(ctx, query,
		token.ID,
		token.CustomerID,
		token.Token,
		token.Platform,
		token.CreatedAt,
		token.UpdatedAt,
	)

	return err
}

func (r *FCMTokenRepositoryImpl) GetByCustomerID(ctx context.Context, customerID string) ([]entity.FCMToken, error) {
	query := `
		SELECT id, customer_id, token, platform, created_at, updated_at
		FROM fcm_tokens
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.Pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []entity.FCMToken
	for rows.Next() {
		var token entity.FCMToken
		if err := rows.Scan(
			&token.ID,
			&token.CustomerID,
			&token.Token,
			&token.Platform,
			&token.CreatedAt,
			&token.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (r *FCMTokenRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM fcm_tokens WHERE id = $1`
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

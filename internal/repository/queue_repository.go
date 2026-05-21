package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueueRepositoryImpl struct {
	Pool *pgxpool.Pool
}

func NewQueueRepository(pool *pgxpool.Pool) *QueueRepositoryImpl {
	return &QueueRepositoryImpl{Pool: pool}
}

func (r *QueueRepositoryImpl) Create(ctx context.Context, queue *entity.Queue) error {
	query := `
		INSERT INTO queues (id, queue_number, customer_id, customer_name, barber_id, service_id, service_name, status, position, created_at, called_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	var barberID *string
	if queue.BarberID != nil {
		barberID = queue.BarberID
	}

	_, err := r.Pool.Exec(ctx, query,
		queue.ID,
		queue.QueueNumber,
		queue.CustomerID,
		queue.CustomerName,
		barberID,
		queue.ServiceID,
		queue.ServiceName,
		queue.Status,
		queue.Position,
		queue.CreatedAt,
		queue.CalledAt,
		queue.CompletedAt,
	)

	return err
}

func (r *QueueRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Queue, error) {
	query := `
		SELECT id, queue_number, customer_id, customer_name, barber_id, service_id, service_name, status, position, created_at, called_at, completed_at
		FROM queues WHERE id = $1
	`

	var queue entity.Queue
	var barberID sql.NullString
	var calledAt sql.NullTime
	var completedAt sql.NullTime

	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&queue.ID,
		&queue.QueueNumber,
		&queue.CustomerID,
		&queue.CustomerName,
		&barberID,
		&queue.ServiceID,
		&queue.ServiceName,
		&queue.Status,
		&queue.Position,
		&queue.CreatedAt,
		&calledAt,
		&completedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("queue not found")
	}
	if err != nil {
		return nil, err
	}

	if barberID.Valid {
		queue.BarberID = &barberID.String
	}
	if calledAt.Valid {
		queue.CalledAt = &calledAt.Time
	}
	if completedAt.Valid {
		queue.CompletedAt = &completedAt.Time
	}

	return &queue, nil
}

func (r *QueueRepositoryImpl) GetList(ctx context.Context) ([]entity.Queue, error) {
	query := `
		SELECT id, queue_number, customer_id, customer_name, barber_id, service_id, service_name, status, position, created_at, called_at, completed_at
		FROM queues
		ORDER BY queue_number ASC
	`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queues []entity.Queue
	for rows.Next() {
		var queue entity.Queue
		var barberID sql.NullString
		var calledAt sql.NullTime
		var completedAt sql.NullTime

		if err := rows.Scan(
			&queue.ID,
			&queue.QueueNumber,
			&queue.CustomerID,
			&queue.CustomerName,
			&barberID,
			&queue.ServiceID,
			&queue.ServiceName,
			&queue.Status,
			&queue.Position,
			&queue.CreatedAt,
			&calledAt,
			&completedAt,
		); err != nil {
			return nil, err
		}

		if barberID.Valid {
			queue.BarberID = &barberID.String
		}
		if calledAt.Valid {
			queue.CalledAt = &calledAt.Time
		}
		if completedAt.Valid {
			queue.CompletedAt = &completedAt.Time
		}

		queues = append(queues, queue)
	}

	return queues, rows.Err()
}

func (r *QueueRepositoryImpl) Update(ctx context.Context, queue *entity.Queue) error {
	query := `
		UPDATE queues
		SET customer_id = $1, customer_name = $2, barber_id = $3, service_id = $4, service_name = $5, status = $6, position = $7, called_at = $8, completed_at = $9
		WHERE id = $10
	`

	var barberID *string
	if queue.BarberID != nil {
		barberID = queue.BarberID
	}

	_, err := r.Pool.Exec(ctx, query,
		queue.CustomerID,
		queue.CustomerName,
		barberID,
		queue.ServiceID,
		queue.ServiceName,
		queue.Status,
		queue.Position,
		queue.CalledAt,
		queue.CompletedAt,
		queue.ID,
	)

	return err
}

func (r *QueueRepositoryImpl) UpdateStatus(ctx context.Context, id string, status entity.QueueStatus) error {
	query := `UPDATE queues SET status = $1 WHERE id = $2`
	_, err := r.Pool.Exec(ctx, query, status, id)
	return err
}

func (r *QueueRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM queues WHERE id = $1`
	_, err := r.Pool.Exec(ctx, query, id)
	return err
}

func (r *QueueRepositoryImpl) GetByCustomerID(ctx context.Context, customerID string) ([]entity.Queue, error) {
	query := `
		SELECT id, queue_number, customer_id, customer_name, barber_id, service_id, service_name, status, position, created_at, called_at, completed_at
		FROM queues WHERE customer_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.Pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queues []entity.Queue
	for rows.Next() {
		var queue entity.Queue
		var barberID sql.NullString
		var calledAt sql.NullTime
		var completedAt sql.NullTime

		if err := rows.Scan(
			&queue.ID,
			&queue.QueueNumber,
			&queue.CustomerID,
			&queue.CustomerName,
			&barberID,
			&queue.ServiceID,
			&queue.ServiceName,
			&queue.Status,
			&queue.Position,
			&queue.CreatedAt,
			&calledAt,
			&completedAt,
		); err != nil {
			return nil, err
		}

		if barberID.Valid {
			queue.BarberID = &barberID.String
		}
		if calledAt.Valid {
			queue.CalledAt = &calledAt.Time
		}
		if completedAt.Valid {
			queue.CompletedAt = &completedAt.Time
		}

		queues = append(queues, queue)
	}

	return queues, rows.Err()
}

func (r *QueueRepositoryImpl) GetCountByStatus(ctx context.Context, status entity.QueueStatus) (int, error) {
	query := `SELECT COUNT(*) FROM queues WHERE status = $1`
	var count int
	err := r.Pool.QueryRow(ctx, query, status).Scan(&count)
	return count, err
}

func (r *QueueRepositoryImpl) GetNextQueueNumber(ctx context.Context) (int, error) {
	query := `SELECT COALESCE(MAX(queue_number), 0) + 1 FROM queues WHERE DATE(created_at) = CURRENT_DATE`
	var nextNumber int
	err := r.Pool.QueryRow(ctx, query).Scan(&nextNumber)
	return nextNumber, err
}

func (r *QueueRepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return nil, fmt.Errorf("transaction not implemented for pgxpool")
}

func (r *QueueRepositoryImpl) GetCompletedToday(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM queues WHERE status = 'completed' AND DATE(completed_at) = CURRENT_DATE`
	var count int
	err := r.Pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *QueueRepositoryImpl) GetCanceledToday(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM queues WHERE status = 'canceled' AND DATE(created_at) = CURRENT_DATE`
	var count int
	err := r.Pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *QueueRepositoryImpl) GetAvgWaitTimeToday(ctx context.Context) (float64, error) {
	query := `
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (called_at - created_at)) / 60), 0)
		FROM queues
		WHERE called_at IS NOT NULL AND DATE(created_at) = CURRENT_DATE
	`
	var avgWait float64
	err := r.Pool.QueryRow(ctx, query).Scan(&avgWait)
	return avgWait, err
}

func (r *QueueRepositoryImpl) GetAvgServiceTimeToday(ctx context.Context) (float64, error) {
	query := `
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - called_at)) / 60), 0)
		FROM queues
		WHERE completed_at IS NOT NULL AND called_at IS NOT NULL AND DATE(created_at) = CURRENT_DATE
	`
	var avgService float64
	err := r.Pool.QueryRow(ctx, query).Scan(&avgService)
	return avgService, err
}

func (r *QueueRepositoryImpl) GetTotalRevenueToday(ctx context.Context) (float64, error) {
	query := `
		SELECT COALESCE(SUM(s.price), 0)
		FROM queues q
		JOIN services s ON q.service_id = s.id
		WHERE q.status = 'completed' AND DATE(q.completed_at) = CURRENT_DATE
	`
	var total float64
	err := r.Pool.QueryRow(ctx, query).Scan(&total)
	return total, err
}

func (r *QueueRepositoryImpl) GetCancelCountByCustomerToday(ctx context.Context, customerID string) (int, error) {
	query := `
		SELECT COUNT(*) FROM queues
		WHERE customer_id = $1 AND status = 'canceled' AND DATE(created_at) = CURRENT_DATE
	`
	var count int
	err := r.Pool.QueryRow(ctx, query, customerID).Scan(&count)
	return count, err
}

func (r *QueueRepositoryImpl) GetLastCancelTimeByCustomer(ctx context.Context, customerID string) (*time.Time, error) {
	query := `
		SELECT created_at FROM queues
		WHERE customer_id = $1 AND status = 'canceled'
		ORDER BY created_at DESC LIMIT 1
	`
	var cancelTime time.Time
	err := r.Pool.QueryRow(ctx, query, customerID).Scan(&cancelTime)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cancelTime, nil
}

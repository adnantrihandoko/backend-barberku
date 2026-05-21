package service

import (
	"context"
	"errors"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/repository"
)

var (
	ErrHistoryNotFound = errors.New("riwayat tidak ditemukan")
)

type HistoryServiceImpl struct {
	queueRepo repository.QueueRepository
}

func NewHistoryService(queueRepo repository.QueueRepository) *HistoryServiceImpl {
	return &HistoryServiceImpl{
		queueRepo: queueRepo,
	}
}

func (s *HistoryServiceImpl) GetHistory(ctx context.Context, customerID string) ([]entity.Queue, error) {
	queues, err := s.queueRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	var history []entity.Queue
	for _, q := range queues {
		if q.Status == string(entity.QueueStatusCompleted) || 
		   q.Status == string(entity.QueueStatusCanceled) || 
		   q.Status == string(entity.QueueStatusSkipped) {
			history = append(history, q)
		}
	}

	return history, nil
}

func (s *HistoryServiceImpl) RateService(ctx context.Context, queueID string, rating int, comment string) error {
	queue, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return ErrHistoryNotFound
	}

	if queue.Status != string(entity.QueueStatusCompleted) {
		return errors.New("hanya bisa memberi rating pada layanan yang selesai")
	}

	if rating < 1 || rating > 5 {
		return errors.New("rating harus antara 1 sampai 5")
	}

	return s.queueRepo.Update(ctx, queue)
}

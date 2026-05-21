package service

import (
	"context"
	"errors"
	"time"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/repository"
)

var (
	ErrQueueNotFound    = errors.New("antrian tidak ditemukan")
	ErrQueueNotWaiting  = errors.New("antrian tidak dalam status menunggu")
	ErrQueueNotCalled   = errors.New("antrian tidak dalam status dipanggil")
	ErrQueueNotServing  = errors.New("antrian tidak dalam status dilayani")
	ErrQueueAlreadyDone = errors.New("antrian sudah selesai")
)

type QueueServiceImpl struct {
	queueRepo repository.QueueRepository
	broadcaster func(event string, data interface{})
}

func NewQueueService(queueRepo repository.QueueRepository, broadcaster func(string, interface{})) *QueueServiceImpl {
	return &QueueServiceImpl{
		queueRepo:   queueRepo,
		broadcaster: broadcaster,
	}
}

func (s *QueueServiceImpl) GetQueueList(ctx context.Context) ([]entity.Queue, error) {
	return s.queueRepo.GetList(ctx)
}

func (s *QueueServiceImpl) GetQueueDetail(ctx context.Context, queueID string) (*entity.Queue, error) {
	return s.queueRepo.GetByID(ctx, queueID)
}

func (s *QueueServiceImpl) JoinQueue(ctx context.Context, customerID, customerName, serviceID, serviceName string, barberID *string) (*entity.Queue, error) {
	nextNumber, err := s.queueRepo.GetNextQueueNumber(ctx)
	if err != nil {
		return nil, err
	}

	queue := &entity.Queue{
		QueueNumber:  nextNumber,
		CustomerID:   customerID,
		CustomerName: customerName,
		ServiceID:    serviceID,
		ServiceName:  serviceName,
		BarberID:     barberID,
		Status:       string(entity.QueueStatusWaiting),
		CreatedAt:    time.Now(),
	}

	if err := s.queueRepo.Create(ctx, queue); err != nil {
		return nil, err
	}

	s.broadcastUpdate()
	return queue, nil
}

func (s *QueueServiceImpl) CallQueue(ctx context.Context, queueID string) error {
	queue, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return ErrQueueNotFound
	}

	if queue.Status != string(entity.QueueStatusWaiting) {
		return ErrQueueNotWaiting
	}

	now := time.Now()
	queue.Status = string(entity.QueueStatusCalled)
	queue.CalledAt = &now

	if err := s.queueRepo.Update(ctx, queue); err != nil {
		return err
	}

	s.broadcastUpdate()
	return nil
}

func (s *QueueServiceImpl) ServeQueue(ctx context.Context, queueID string) error {
	queue, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return ErrQueueNotFound
	}

	if queue.Status != string(entity.QueueStatusCalled) {
		return ErrQueueNotCalled
	}

	queue.Status = string(entity.QueueStatusInProgress)

	if err := s.queueRepo.Update(ctx, queue); err != nil {
		return err
	}

	s.broadcastUpdate()
	return nil
}

func (s *QueueServiceImpl) CompleteQueue(ctx context.Context, queueID string) error {
	queue, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return ErrQueueNotFound
	}

	if queue.Status != string(entity.QueueStatusInProgress) && queue.Status != string(entity.QueueStatusCalled) {
		return ErrQueueNotServing
	}

	now := time.Now()
	queue.Status = string(entity.QueueStatusCompleted)
	queue.CompletedAt = &now

	if err := s.queueRepo.Update(ctx, queue); err != nil {
		return err
	}

	s.broadcastUpdate()
	return nil
}

func (s *QueueServiceImpl) SkipQueue(ctx context.Context, queueID string) error {
	queue, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return ErrQueueNotFound
	}

	queue.Status = string(entity.QueueStatusSkipped)

	if err := s.queueRepo.Update(ctx, queue); err != nil {
		return err
	}

	s.broadcastUpdate()
	return nil
}

func (s *QueueServiceImpl) CancelQueue(ctx context.Context, queueID string) error {
	queue, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return ErrQueueNotFound
	}

	if queue.Status == string(entity.QueueStatusCompleted) {
		return ErrQueueAlreadyDone
	}

	queue.Status = string(entity.QueueStatusCanceled)

	if err := s.queueRepo.Update(ctx, queue); err != nil {
		return err
	}

	s.broadcastUpdate()
	return nil
}

func (s *QueueServiceImpl) AddWalkIn(ctx context.Context, customerName, serviceID, serviceName string, barberID *string) (*entity.Queue, error) {
	nextNumber, err := s.queueRepo.GetNextQueueNumber(ctx)
	if err != nil {
		return nil, err
	}

	queue := &entity.Queue{
		QueueNumber:  nextNumber,
		CustomerID:   "walk-in",
		CustomerName: customerName,
		ServiceID:    serviceID,
		ServiceName:  serviceName,
		BarberID:     barberID,
		Status:       string(entity.QueueStatusWaiting),
		CreatedAt:    time.Now(),
	}

	if err := s.queueRepo.Create(ctx, queue); err != nil {
		return nil, err
	}

	s.broadcastUpdate()
	return queue, nil
}

func (s *QueueServiceImpl) broadcastUpdate() {
	if s.broadcaster != nil {
		go func() {
			s.broadcaster("queue_updated", nil)
		}()
	}
}

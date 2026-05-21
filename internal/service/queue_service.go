package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/repository"
)

var (
	ErrQueueNotFound       = errors.New("antrian tidak ditemukan")
	ErrQueueNotWaiting     = errors.New("antrian tidak dalam status menunggu")
	ErrQueueNotCalled      = errors.New("antrian tidak dalam status dipanggil")
	ErrQueueNotServing     = errors.New("antrian tidak dalam status dilayani")
	ErrQueueAlreadyDone    = errors.New("antrian sudah selesai")
	ErrOutsideBusinessHour = errors.New("di luar jam operasional")
	ErrQueueFull           = errors.New("antrian sudah penuh")
	ErrActiveQueueExists   = errors.New("anda sudah memiliki antrian aktif")
	ErrCancelLimitReached  = errors.New("batas pembatalan harian tercapai. Silakan hubungi barber")
	ErrCancelCooldown      = errors.New("tidak bisa bergabung kembali, harap tunggu 15 menit setelah pembatalan")
)

type QueueServiceImpl struct {
	queueRepo     repository.QueueRepository
	settingsRepo  repository.StoreSettingsRepository
	broadcaster   func(event string, data interface{})
	fcmService    FCMService
}

func NewQueueService(
	queueRepo repository.QueueRepository,
	settingsRepo repository.StoreSettingsRepository,
	broadcaster func(string, interface{}),
	fcmService FCMService,
) *QueueServiceImpl {
	return &QueueServiceImpl{
		queueRepo:    queueRepo,
		settingsRepo: settingsRepo,
		broadcaster:  broadcaster,
		fcmService:   fcmService,
	}
}

func (s *QueueServiceImpl) GetQueueList(ctx context.Context) ([]entity.Queue, error) {
	return s.queueRepo.GetList(ctx)
}

func (s *QueueServiceImpl) GetQueueDetail(ctx context.Context, queueID string) (*entity.Queue, error) {
	return s.queueRepo.GetByID(ctx, queueID)
}

func (s *QueueServiceImpl) JoinQueue(ctx context.Context, customerID, customerName, serviceID, serviceName string, barberID *string) (*entity.Queue, error) {
	now := time.Now()

	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		settings = &entity.StoreSettings{
			OpenHour:     9,
			CloseHour:    21,
			MaxQueueSize: 50,
		}
	}

	if now.Hour() < settings.OpenHour || now.Hour() >= settings.CloseHour {
		return nil, ErrOutsideBusinessHour
	}

	waitingCount, err := s.queueRepo.GetCountByStatus(ctx, entity.QueueStatusWaiting)
	if err != nil {
		return nil, err
	}
	if waitingCount >= settings.MaxQueueSize {
		return nil, ErrQueueFull
	}

	activeQueues, err := s.queueRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	for _, q := range activeQueues {
		if q.Status == string(entity.QueueStatusWaiting) || q.Status == string(entity.QueueStatusCalled) || q.Status == string(entity.QueueStatusInProgress) {
			return nil, ErrActiveQueueExists
		}
	}

	cancelCount, err := s.queueRepo.GetCancelCountByCustomerToday(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if cancelCount >= 2 {
		return nil, ErrCancelLimitReached
	}

	lastCancelTime, err := s.queueRepo.GetLastCancelTimeByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if lastCancelTime != nil && now.Sub(*lastCancelTime) < 15*time.Minute {
		return nil, ErrCancelCooldown
	}

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
		CreatedAt:    now,
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

	if s.fcmService != nil {
		go func() {
			title := "Giliran Anda Tiba!"
			body := fmt.Sprintf("Nomor antrian %d, silakan ke kasir.", queue.QueueNumber)
			if err := s.fcmService.SendNotification(context.Background(), queue.CustomerID, title, body); err != nil {
				slog.Error("failed to send FCM notification", "queue_id", queueID, "error", err)
			}
		}()
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

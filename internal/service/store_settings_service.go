package service

import (
	"context"
	"errors"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/repository"
)

var (
	ErrSettingsNotFound = errors.New("pengaturan toko tidak ditemukan")
)

type StoreSettingsServiceImpl struct {
	settingsRepo repository.StoreSettingsRepository
}

func NewStoreSettingsService(settingsRepo repository.StoreSettingsRepository) *StoreSettingsServiceImpl {
	return &StoreSettingsServiceImpl{
		settingsRepo: settingsRepo,
	}
}

func (s *StoreSettingsServiceImpl) GetSettings(ctx context.Context) (*entity.StoreSettings, error) {
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return &entity.StoreSettings{
			OpenHour:     9,
			CloseHour:    21,
			MaxQueueSize: 50,
		}, nil
	}
	return settings, nil
}

func (s *StoreSettingsServiceImpl) UpdateSettings(ctx context.Context, settings *entity.StoreSettings) error {
	if settings.OpenHour < 0 || settings.OpenHour > 23 || settings.CloseHour < 0 || settings.CloseHour > 23 {
		return errors.New("jam operasional tidak valid")
	}
	if settings.OpenHour >= settings.CloseHour {
		return errors.New("jam buka harus lebih awal dari jam tutup")
	}
	if settings.MaxQueueSize < 1 {
		return errors.New("batas antrian minimal 1")
	}

	return s.settingsRepo.Update(ctx, settings)
}

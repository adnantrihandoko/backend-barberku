package service

import (
	"context"
	"errors"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/repository"
)

var (
	ErrBarberNotFound = errors.New("barber tidak ditemukan")
)

type BarberServiceImpl struct {
	barberRepo repository.BarberRepository
}

func NewBarberService(barberRepo repository.BarberRepository) *BarberServiceImpl {
	return &BarberServiceImpl{
		barberRepo: barberRepo,
	}
}

func (s *BarberServiceImpl) ListBarbers(ctx context.Context) ([]entity.Barber, error) {
	return s.barberRepo.List(ctx)
}

func (s *BarberServiceImpl) GetBarber(ctx context.Context, id string) (*entity.Barber, error) {
	return s.barberRepo.GetByID(ctx, id)
}

func (s *BarberServiceImpl) CreateBarber(ctx context.Context, name, specialty string) (*entity.Barber, error) {
	barber := &entity.Barber{
		Name:      name,
		Specialty: specialty,
		IsActive:  true,
	}

	if err := s.barberRepo.Create(ctx, barber); err != nil {
		return nil, err
	}

	return barber, nil
}

func (s *BarberServiceImpl) UpdateBarber(ctx context.Context, id, name, specialty string, isActive bool) (*entity.Barber, error) {
	barber, err := s.barberRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrBarberNotFound
	}

	barber.Name = name
	barber.Specialty = specialty
	barber.IsActive = isActive

	if err := s.barberRepo.Update(ctx, barber); err != nil {
		return nil, err
	}

	return barber, nil
}

func (s *BarberServiceImpl) DeleteBarber(ctx context.Context, id string) error {
	_, err := s.barberRepo.GetByID(ctx, id)
	if err != nil {
		return ErrBarberNotFound
	}

	return s.barberRepo.Delete(ctx, id)
}

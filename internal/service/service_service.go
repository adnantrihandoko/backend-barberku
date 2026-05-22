package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/repository"
)

var (
	ErrServiceNotFound = errors.New("layanan tidak ditemukan")
)

type ServiceServiceImpl struct {
	serviceRepo repository.ServiceRepository
}

func NewServiceService(serviceRepo repository.ServiceRepository) *ServiceServiceImpl {
	return &ServiceServiceImpl{
		serviceRepo: serviceRepo,
	}
}

func (s *ServiceServiceImpl) ListServices(ctx context.Context) ([]entity.Service, error) {
	return s.serviceRepo.List(ctx)
}

func (s *ServiceServiceImpl) GetService(ctx context.Context, id string) (*entity.Service, error) {
	return s.serviceRepo.GetByID(ctx, id)
}

func (s *ServiceServiceImpl) CreateService(ctx context.Context, name, description string, price float64, duration int) (*entity.Service, error) {
	now := time.Now()
	service := &entity.Service{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		Duration:    duration,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.serviceRepo.Create(ctx, service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceServiceImpl) UpdateService(ctx context.Context, id, name, description string, price float64, duration int, isActive bool) (*entity.Service, error) {
	service, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrServiceNotFound
	}

	service.Name = name
	service.Description = description
	service.Price = price
	service.Duration = duration
	service.IsActive = isActive

	if err := s.serviceRepo.Update(ctx, service); err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServiceServiceImpl) DeleteService(ctx context.Context, id string) error {
	_, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		return ErrServiceNotFound
	}

	return s.serviceRepo.Delete(ctx, id)
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type DeviceService struct {
	deviceRepo domain.DeviceRepository
}

func NewDeviceService(conn repository.Connection) domain.DeviceService {
	return &DeviceService{
		deviceRepo: repository.NewPostgresDeviceRepository(conn),
	}
}

func (s *DeviceService) GetByID(ctx context.Context, userID, deviceID int64) (*domain.Device, error) {
	if userID < 0 {
		return nil, fmt.Errorf("%w: user ID cannot be negative", domain.ErrValidation)
	}

	if deviceID < 0 {
		return nil, fmt.Errorf("%w: device ID cannot be negative", domain.ErrValidation)
	}

	return s.deviceRepo.GetByID(ctx, userID, deviceID)
}

func (s *DeviceService) List(ctx context.Context, userID int64) ([]*domain.Device, error) {
	if userID < 0 {
		return nil, fmt.Errorf("%w: user ID cannot be negative", domain.ErrValidation)
	}

	return s.deviceRepo.List(ctx, userID)
}

func (s *DeviceService) Create(ctx context.Context, userID int64, name, platform, model string) error {
	now := time.Now().UTC()

	device := &domain.Device{
		UserID:    userID,
		Name:      name,
		Platform:  domain.Platform(platform),
		Model:     model,
		Active:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := device.Validate(); err != nil {
		return err
	}

	return s.deviceRepo.Create(ctx, device)
}

func (s *DeviceService) Update(ctx context.Context, userID, deviceID int64, name, token string, acive bool) error {

	device, err := s.deviceRepo.GetByID(ctx, userID, deviceID)
	if err != nil {
		return err
	}

	if device.Name != name {
		device.Name = name
	}

	if device.Token != &token {
		device.Token = &token
	}

	if device.Active != acive {
		device.Active = acive
	}

	if err := device.Validate(); err != nil {
		return err
	}

	return s.deviceRepo.Update(ctx, device)
}

func (s *DeviceService) Delete(ctx context.Context, userID, deviceID int64) error {
	if userID < 0 {
		return fmt.Errorf("%w: user ID cannot be negative", domain.ErrValidation)
	}

	if deviceID < 0 {
		return fmt.Errorf("%w: device ID cannot be negative", domain.ErrValidation)
	}

	return s.deviceRepo.Delete(ctx, userID, deviceID)
}

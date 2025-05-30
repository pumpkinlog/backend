package service

import (
	"context"
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

func (s *DeviceService) Create(ctx context.Context, userID, name, platform, model string) error {
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

func (s *DeviceService) Update(ctx context.Context, userID, deviceID, name, token string, acive bool) error {

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

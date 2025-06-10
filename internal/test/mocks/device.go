package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type DeviceRepository struct {
	GetByIDFunc func(ctx context.Context, userID, deviceID int64) (*domain.Device, error)
	ListFunc    func(ctx context.Context, userID int64) ([]*domain.Device, error)
	CreateFunc  func(ctx context.Context, device *domain.Device) error
	UpdateFunc  func(ctx context.Context, device *domain.Device) error
	DeleteFunc  func(ctx context.Context, userID, deviceID int64) error
}

func (m DeviceRepository) GetByID(ctx context.Context, userID, deviceID int64) (*domain.Device, error) {
	return m.GetByIDFunc(ctx, userID, deviceID)
}

func (m DeviceRepository) List(ctx context.Context, userID int64) ([]*domain.Device, error) {
	return m.ListFunc(ctx, userID)
}

func (m DeviceRepository) Create(ctx context.Context, device *domain.Device) error {
	return m.CreateFunc(ctx, device)
}

func (m DeviceRepository) Update(ctx context.Context, device *domain.Device) error {
	return m.UpdateFunc(ctx, device)
}

func (m DeviceRepository) Delete(ctx context.Context, userID, deviceID int64) error {
	return m.DeleteFunc(ctx, userID, deviceID)
}

type DeviceService struct {
	GetByIDFunc func(ctx context.Context, userID, deviceID int64) (*domain.Device, error)
	ListFunc    func(ctx context.Context, userID int64) ([]*domain.Device, error)
	CreateFunc  func(ctx context.Context, userID int64, name, platform, model string) error
	UpdateFunc  func(ctx context.Context, userID, deviceID int64, name, token string, acive bool) error
	DeleteFunc  func(ctx context.Context, userID, deviceID int64) error
}

func (m DeviceService) GetByID(ctx context.Context, userID, deviceID int64) (*domain.Device, error) {
	return m.GetByIDFunc(ctx, userID, deviceID)
}

func (m DeviceService) List(ctx context.Context, userID int64) ([]*domain.Device, error) {
	return m.ListFunc(ctx, userID)
}

func (m DeviceService) Create(ctx context.Context, userID int64, name, platform, model string) error {
	return m.CreateFunc(ctx, userID, name, platform, model)
}

func (m DeviceService) Update(ctx context.Context, userID, deviceID int64, name, token string, acive bool) error {
	return m.UpdateFunc(ctx, userID, deviceID, name, token, acive)
}

func (m DeviceService) Delete(ctx context.Context, userID, deviceID int64) error {
	return m.DeleteFunc(ctx, userID, deviceID)
}

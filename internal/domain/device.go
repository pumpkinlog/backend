package domain

import (
	"context"
	"fmt"
	"time"
)

type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
)

type Device struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	Name      string    `json:"name"`
	Platform  Platform  `json:"platform"`
	Model     string    `json:"model"`
	Token     *string   `json:"token"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p Platform) Valid() bool {
	switch p {
	case PlatformIOS, PlatformAndroid:
		return true
	default:
		return false
	}
}

func (d *Device) Validate() error {

	if d.UserID < 0 {
		return fmt.Errorf("%w: user ID is invalid", ErrValidation)
	}

	if d.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}

	if !d.Platform.Valid() {
		return fmt.Errorf("%w: platform is invalid", ErrValidation)
	}

	if d.Model == "" {
		return fmt.Errorf("%w: model is required", ErrValidation)
	}

	if d.CreatedAt.IsZero() {
		return fmt.Errorf("%w: created at timestamp is invalid", ErrValidation)
	}

	if d.UpdatedAt.IsZero() {
		return fmt.Errorf("%w: updated at timestamp is invalid", ErrValidation)
	}

	now := time.Now().UTC()

	if d.CreatedAt.After(now) {
		return fmt.Errorf("%w: created at timestamp cannot be in the future", ErrValidation)
	}

	if d.UpdatedAt.After(now) {
		return fmt.Errorf("%w: updated at timestamp cannot be in the future", ErrValidation)
	}

	return nil
}

type DeviceService interface {
	GetByID(ctx context.Context, userID, deviceID int64) (*Device, error)
	List(ctx context.Context, userID int64) ([]*Device, error)
	Create(ctx context.Context, userID int64, name, platform, model string) error
	Update(ctx context.Context, userID, deviceID int64, name, token string, acive bool) error
	Delete(ctx context.Context, userID, deviceID int64) error
}

type DeviceRepository interface {
	GetByID(ctx context.Context, userID, deviceID int64) (*Device, error)
	List(ctx context.Context, userID int64) ([]*Device, error)
	Create(ctx context.Context, device *Device) error
	Update(ctx context.Context, device *Device) error
	Delete(ctx context.Context, userID, deviceID int64) error
}

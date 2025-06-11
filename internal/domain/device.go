package domain

import (
	"context"
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
	if d.UserID <= 0 {
		return ValidationError("user ID is required")
	}

	if d.Name == "" {
		return ValidationError("name is required")
	}

	if !d.Platform.Valid() {
		return ValidationError("platform is invalid")
	}

	if d.Model == "" {
		return ValidationError("model is required")
	}

	if d.CreatedAt.IsZero() {
		return ValidationError("created at timestamp is invalid")
	}

	if d.UpdatedAt.IsZero() {
		return ValidationError("updated at timestamp is invalid")
	}

	now := time.Now().UTC()

	if d.CreatedAt.After(now) {
		return ValidationError("created at timestamp cannot be in the future")
	}

	if d.UpdatedAt.After(now) {
		return ValidationError("updated at timestamp cannot be in the future")
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

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
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Platform  Platform  `json:"platform"`
	Model     string    `json:"model"`
	Token     *string   `json:"token"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeviceRepository interface {
	GetByID(ctx context.Context, userID, deviceID string) (*Device, error)
	List(ctx context.Context, userID string) ([]*Device, error)

	Create(ctx context.Context, device *Device) error
	Update(ctx context.Context, device *Device) error

	Delete(ctx context.Context, userID, deviceID string) error
}

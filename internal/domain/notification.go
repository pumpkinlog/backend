package domain

import (
	"context"
	"time"
)

type Notification struct {
	ID       string    `json:"id"`
	DeviceID string    `json:"device_id"`
	Message  string    `json:"message"`
	Body     string    `json:"body"`
	Date     time.Time `json:"date"`
}

type NotificationRepository interface {
	GetByID(ctx context.Context, id string) (*Notification, error)

	List(ctx context.Context, deviceID string) ([]*Notification, error)

	Create(ctx context.Context, notification *Notification) error
	Update(ctx context.Context, notification *Notification) error

	Delete(ctx context.Context, id string) error
}

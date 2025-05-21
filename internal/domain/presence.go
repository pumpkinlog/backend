package domain

import (
	"context"
	"time"
)

type Presence struct {
	UserID   string    `json:"userId"`
	RegionID string    `json:"regionId"`
	Date     time.Time `json:"date"`
	DeviceID *string   `json:"deviceId,omitempty"`
}

type PresenceFilter struct {
	RegionIDs []string
	Start     *time.Time
	End       *time.Time
	Page      *int
	Limit     *int
}

type PresenceService interface {
	Create(ctx context.Context, userID, regionID string, deviceID *string, start, end time.Time) error
	Delete(ctx context.Context, userID, regionID string, start, end time.Time) error
}

type PresenceRepository interface {
	GetByID(ctx context.Context, userID, regionID string, date time.Time) (*Presence, error)
	List(ctx context.Context, userID string, filter *PresenceFilter) ([]*Presence, error)
	ListByRegionBounds(ctx context.Context, userID string, bounds map[string]TimeWindow) ([]*Presence, error)

	Create(ctx context.Context, location *Presence) error
	CreateRange(ctx context.Context, userID, regionID string, deviceID *string, start, end time.Time) error

	Delete(ctx context.Context, userID, regionID string, date time.Time) error
	DeleteRange(ctx context.Context, userID, regionID string, start, end time.Time) error
}

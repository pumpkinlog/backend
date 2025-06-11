package domain

import (
	"context"
	"time"
)

type Presence struct {
	UserID    int64     `json:"userId"`
	RegionID  RegionID  `json:"regionId"`
	Date      time.Time `json:"date"`
	DeviceID  *string   `json:"deviceId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p *Presence) Validate() error {

	if p.UserID < 0 {
		return ValidationError("user ID is required")
	}

	if err := p.RegionID.Validate(); err != nil {
		return err
	}

	if p.Date.After(time.Now().UTC()) {
		return ValidationError("date cannot be in the future")
	}

	if p.DeviceID != nil && *p.DeviceID == "" {
		return ValidationError("device ID cannot be empty")
	}

	if p.CreatedAt.IsZero() {
		return ValidationError("created at timestamp is required")
	}

	if p.UpdatedAt.IsZero() {
		return ValidationError("updated at timestamp is required")
	}

	now := time.Now().UTC()

	if p.CreatedAt.After(now) {
		return ValidationError("created at timestamp cannot be in the future")
	}

	if p.UpdatedAt.After(now) {
		return ValidationError("updated at timestamp cannot be in the future")
	}

	return nil
}

type PresenceService interface {
	GetByID(ctx context.Context, userID int64, regionID RegionID, date time.Time) (*Presence, error)
	List(ctx context.Context, userID int64, filter *PresenceFilter) ([]*Presence, error)
	Create(ctx context.Context, userID int64, regionID RegionID, deviceID *int64, start, end time.Time) error
	Delete(ctx context.Context, userID int64, regionID RegionID, start, end time.Time) error
}

type PresenceFilter struct {
	RegionIDs []RegionID
	Start     *time.Time
	End       *time.Time
}

type PresenceRepository interface {
	GetByID(ctx context.Context, userID int64, regionID RegionID, date time.Time) (*Presence, error)
	List(ctx context.Context, userID int64, filter *PresenceFilter) ([]*Presence, error)
	ListByRegionPeriod(ctx context.Context, userID int64, regionID RegionID, start, end time.Time) ([]*Presence, error)
	Create(ctx context.Context, location *Presence) error
	CreateRange(ctx context.Context, userID int64, regionID RegionID, deviceID *int64, start, end time.Time) error
	Delete(ctx context.Context, userID int64, regionID RegionID, date time.Time) error
	DeleteRange(ctx context.Context, userID int64, regionID RegionID, start, end time.Time) error
}

package domain

import (
	"context"
	"fmt"
	"time"
)

type Presence struct {
	UserID    string    `json:"userId"`
	RegionID  string    `json:"regionId"`
	Date      time.Time `json:"date"`
	DeviceID  *string   `json:"deviceId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (p *Presence) Validate() error {

	if p.UserID == "" {
		return fmt.Errorf("%w: user ID is required", ErrValidation)
	}

	length := len(p.RegionID)
	if length < 3 || length > 5 {
		return fmt.Errorf("%w: region ID must be between 3-5 characters", ErrValidation)
	}

	if p.Date.After(time.Now().UTC()) {
		return fmt.Errorf("%w: date cannot be in the future", ErrValidation)
	}

	if p.DeviceID != nil && *p.DeviceID == "" {
		return fmt.Errorf("%w: device ID cannot be an empty string", ErrValidation)
	}

	if p.CreatedAt.IsZero() {
		return fmt.Errorf("%w: created at timestamp is invalid", ErrValidation)
	}

	if p.UpdatedAt.IsZero() {
		return fmt.Errorf("%w: updated at timestamp is invalid", ErrValidation)
	}

	now := time.Now().UTC()

	if p.CreatedAt.After(now) {
		return fmt.Errorf("%w: created at timestamp cannot be in the future", ErrValidation)
	}

	if p.UpdatedAt.After(now) {
		return fmt.Errorf("%w: updated at timestamp cannot be in the future", ErrValidation)
	}

	return nil
}

type PresenceService interface {
	Create(ctx context.Context, userID, regionID string, deviceID *string, start, end time.Time) error
	Delete(ctx context.Context, userID, regionID string, start, end time.Time) error
}

type PresenceFilter struct {
	RegionIDs []string
	Start     *time.Time
	End       *time.Time
	Page      *int
	Limit     *int
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

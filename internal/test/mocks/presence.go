package mocks

import (
	"context"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

type PresenceRepo struct {
	GetByIDFunc            func(ctx context.Context, userID, regionID string, date time.Time) (*domain.Presence, error)
	ListFunc               func(ctx context.Context, id string, filter *domain.PresenceFilter) ([]*domain.Presence, error)
	ListByRegionBoundsFunc func(ctx context.Context, userID string, bounds map[string]domain.TimeWindow) ([]*domain.Presence, error)
	CreateFunc             func(ctx context.Context, location *domain.Presence) error
	CreateRangeFunc        func(ctx context.Context, userID, regionID string, deviceID *string, start, end time.Time) error
	DeleteFunc             func(ctx context.Context, userID, regionID string, date time.Time) error
	DeleteRangeFunc        func(ctx context.Context, userID, regionID string, start, end time.Time) error
}

func (m PresenceRepo) GetByID(ctx context.Context, userID, regionID string, date time.Time) (*domain.Presence, error) {
	return m.GetByIDFunc(ctx, userID, regionID, date)
}

func (m PresenceRepo) List(ctx context.Context, id string, filter *domain.PresenceFilter) ([]*domain.Presence, error) {
	return m.ListFunc(ctx, id, filter)
}

func (m PresenceRepo) ListByRegionBounds(ctx context.Context, userID string, bounds map[string]domain.TimeWindow) ([]*domain.Presence, error) {
	return m.ListByRegionBoundsFunc(ctx, userID, bounds)
}

func (m PresenceRepo) Create(ctx context.Context, location *domain.Presence) error {
	return m.CreateFunc(ctx, location)
}

func (m PresenceRepo) CreateRange(ctx context.Context, userID, regionID string, deviceID *string, start, end time.Time) error {
	return m.CreateRangeFunc(ctx, userID, regionID, deviceID, start, end)
}

func (m PresenceRepo) Delete(ctx context.Context, userID, regionID string, date time.Time) error {
	return m.DeleteFunc(ctx, userID, regionID, date)
}

func (m PresenceRepo) DeleteRange(ctx context.Context, userID, regionID string, start, end time.Time) error {
	return m.DeleteRangeFunc(ctx, userID, regionID, start, end)
}

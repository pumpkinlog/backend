package mocks

import (
	"context"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

type PresenceRepo struct {
	GetByIDFunc            func(ctx context.Context, userID int64, regionID domain.RegionID, date time.Time) (*domain.Presence, error)
	ListFunc               func(ctx context.Context, userID int64, filter *domain.PresenceFilter) ([]*domain.Presence, error)
	ListByRegionPeriodFunc func(ctx context.Context, userID int64, regionID domain.RegionID, start, end time.Time) ([]*domain.Presence, error)
	CreateFunc             func(ctx context.Context, location *domain.Presence) error
	CreateRangeFunc        func(ctx context.Context, userID int64, regionID domain.RegionID, deviceID *int64, start, end time.Time) error
	DeleteFunc             func(ctx context.Context, userID int64, regionID domain.RegionID, date time.Time) error
	DeleteRangeFunc        func(ctx context.Context, userID int64, regionID domain.RegionID, start, end time.Time) error
}

func (m PresenceRepo) GetByID(ctx context.Context, userID int64, regionID domain.RegionID, date time.Time) (*domain.Presence, error) {
	return m.GetByIDFunc(ctx, userID, regionID, date)
}

func (m PresenceRepo) List(ctx context.Context, userID int64, filter *domain.PresenceFilter) ([]*domain.Presence, error) {
	return m.ListFunc(ctx, userID, filter)
}

func (m PresenceRepo) ListByRegionPeriod(ctx context.Context, userID int64, regionID domain.RegionID, start, end time.Time) ([]*domain.Presence, error) {
	return m.ListByRegionPeriodFunc(ctx, userID, regionID, start, end)
}

func (m PresenceRepo) Create(ctx context.Context, location *domain.Presence) error {
	return m.CreateFunc(ctx, location)
}

func (m PresenceRepo) CreateRange(ctx context.Context, userID int64, regionID domain.RegionID, deviceID *int64, start, end time.Time) error {
	return m.CreateRangeFunc(ctx, userID, regionID, deviceID, start, end)
}

func (m PresenceRepo) Delete(ctx context.Context, userID int64, regionID domain.RegionID, date time.Time) error {
	return m.DeleteFunc(ctx, userID, regionID, date)
}

func (m PresenceRepo) DeleteRange(ctx context.Context, userID int64, regionID domain.RegionID, start, end time.Time) error {
	return m.DeleteRangeFunc(ctx, userID, regionID, start, end)
}

type PresenceService struct {
	GetByIDFunc func(ctx context.Context, userID int64, regionID domain.RegionID, date time.Time) (*domain.Presence, error)
	ListFunc    func(ctx context.Context, userID int64, filter *domain.PresenceFilter) ([]*domain.Presence, error)
	CreateFunc  func(ctx context.Context, userID int64, regionID domain.RegionID, deviceID *int64, start, end time.Time) error
	DeleteFunc  func(ctx context.Context, userID int64, regionID domain.RegionID, start, end time.Time) error
}

func (m PresenceService) GetByID(ctx context.Context, userID int64, regionID domain.RegionID, date time.Time) (*domain.Presence, error) {
	return m.GetByIDFunc(ctx, userID, regionID, date)
}

func (m PresenceService) List(ctx context.Context, userID int64, filter *domain.PresenceFilter) ([]*domain.Presence, error) {
	return m.ListFunc(ctx, userID, filter)
}

func (m PresenceService) Create(ctx context.Context, userID int64, regionID domain.RegionID, deviceID *int64, start, end time.Time) error {
	return m.CreateFunc(ctx, userID, regionID, deviceID, start, end)
}

func (m PresenceService) Delete(ctx context.Context, userID int64, regionID domain.RegionID, start, end time.Time) error {
	return m.DeleteFunc(ctx, userID, regionID, start, end)
}

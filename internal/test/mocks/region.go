package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type RegionRepo struct {
	GetByIDFunc        func(ctx context.Context, regionID domain.RegionID) (*domain.Region, error)
	ListFunc           func(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error)
	CreateOrUpdateFunc func(ctx context.Context, region *domain.Region) error
}

func (m RegionRepo) GetByID(ctx context.Context, id domain.RegionID) (*domain.Region, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m RegionRepo) List(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error) {
	return m.ListFunc(ctx, filter)
}

func (m RegionRepo) CreateOrUpdate(ctx context.Context, region *domain.Region) error {
	return m.CreateOrUpdateFunc(ctx, region)
}

type RegionService struct {
	GetByIDFunc        func(ctx context.Context, id domain.RegionID) (*domain.Region, error)
	ListFunc           func(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error)
	CreateOrUpdateFunc func(ctx context.Context, region *domain.Region) error
}

func (m RegionService) GetByID(ctx context.Context, id domain.RegionID) (*domain.Region, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m RegionService) List(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error) {
	return m.ListFunc(ctx, filter)
}

func (m RegionService) CreateOrUpdate(ctx context.Context, region *domain.Region) error {
	return m.CreateOrUpdateFunc(ctx, region)
}

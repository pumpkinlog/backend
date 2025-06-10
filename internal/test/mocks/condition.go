package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type ConditionRepository struct {
	GetByIDFunc        func(ctx context.Context, id domain.Code) (*domain.Condition, error)
	ListFunc           func(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error)
	ListByRegionIDFunc func(ctx context.Context, regionID string) ([]*domain.Condition, error)
	CreateOrUpdateFunc func(ctx context.Context, condition *domain.Condition) error
	DeleteFunc         func(ctx context.Context, id domain.Code) error
}

func (m ConditionRepository) GetByID(ctx context.Context, id domain.Code) (*domain.Condition, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m ConditionRepository) List(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error) {
	return m.ListFunc(ctx, filter)
}

func (m ConditionRepository) ListByRegionID(ctx context.Context, regionID string) ([]*domain.Condition, error) {
	return m.ListByRegionIDFunc(ctx, regionID)
}

func (m ConditionRepository) CreateOrUpdate(ctx context.Context, condition *domain.Condition) error {
	return m.CreateOrUpdateFunc(ctx, condition)
}

func (m ConditionRepository) Delete(ctx context.Context, id domain.Code) error {
	return m.DeleteFunc(ctx, id)
}

type ConditionService struct {
	GetByIDFunc        func(ctx context.Context, id domain.Code) (*domain.Condition, error)
	ListFunc           func(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error)
	CreateOrUpdateFunc func(ctx context.Context, condition *domain.Condition) error
}

func (m ConditionService) GetByID(ctx context.Context, id domain.Code) (*domain.Condition, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m ConditionService) List(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error) {
	return m.ListFunc(ctx, filter)
}

func (m ConditionService) CreateOrUpdate(ctx context.Context, condition *domain.Condition) error {
	return m.CreateOrUpdateFunc(ctx, condition)
}

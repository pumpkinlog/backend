package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type ConditionRepo struct {
	GetByIDFunc        func(ctx context.Context, id string) (*domain.Condition, error)
	ListFunc           func(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error)
	CreateOrUpdateFunc func(ctx context.Context, condition *domain.Condition) error
	DeleteFunc         func(ctx context.Context, id string) error
}

func (m ConditionRepo) GetByID(ctx context.Context, id string) (*domain.Condition, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m ConditionRepo) List(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error) {
	return m.ListFunc(ctx, filter)
}

func (m ConditionRepo) CreateOrUpdate(ctx context.Context, condition *domain.Condition) error {
	return m.CreateOrUpdateFunc(ctx, condition)
}

func (m ConditionRepo) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

type ConditionService struct {
	CreateFunc func(ctx context.Context, condition *domain.Condition, ruleIDs []string) error
	DeleteFunc func(ctx context.Context, conditionID, ruleID string) error
}

func (m ConditionService) Create(ctx context.Context, condition *domain.Condition, ruleIDs []string) error {
	return m.CreateFunc(ctx, condition, ruleIDs)
}

func (m ConditionService) Delete(ctx context.Context, conditionID, ruleID string) error {
	return m.DeleteFunc(ctx, conditionID, ruleID)
}

package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type RuleRepo struct {
	GetByIDFunc        func(ctx context.Context, ruleID domain.Code) (*domain.Rule, error)
	ListFunc           func(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error)
	ListByRegionIDFunc func(ctx context.Context, regionID string) ([]*domain.Rule, error)
	CreateOrUpdateFunc func(ctx context.Context, rule *domain.Rule) error
}

func (m RuleRepo) GetByID(ctx context.Context, ruleID domain.Code) (*domain.Rule, error) {
	return m.GetByIDFunc(ctx, ruleID)
}

func (m RuleRepo) List(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error) {
	return m.ListFunc(ctx, filter)
}

func (m RuleRepo) ListByRegionID(ctx context.Context, regionID string) ([]*domain.Rule, error) {
	return m.ListByRegionIDFunc(ctx, regionID)
}

func (m RuleRepo) CreateOrUpdate(ctx context.Context, rule *domain.Rule) error {
	return m.CreateOrUpdateFunc(ctx, rule)
}

type RuleService struct {
	GetByIDFunc        func(ctx context.Context, ruleID domain.Code) (*domain.Rule, error)
	ListFunc           func(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error)
	CreateOrUpdateFunc func(ctx context.Context, rule *domain.Rule) error
}

func (m RuleService) GetByID(ctx context.Context, ruleID domain.Code) (*domain.Rule, error) {
	return m.GetByIDFunc(ctx, ruleID)
}

func (m RuleService) List(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error) {
	return m.ListFunc(ctx, filter)
}

func (m RuleService) CreateOrUpdate(ctx context.Context, rule *domain.Rule) error {
	return m.CreateOrUpdateFunc(ctx, rule)
}

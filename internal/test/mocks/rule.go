package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type RuleRepo struct {
	GetByIDFunc        func(ctx context.Context, id string) (*domain.Rule, error)
	ListFunc           func(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error)
	CreateOrUpdateFunc func(ctx context.Context, rule *domain.Rule) error
}

func (m RuleRepo) GetByID(ctx context.Context, id string) (*domain.Rule, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m RuleRepo) List(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error) {
	return m.ListFunc(ctx, filter)
}

func (m RuleRepo) CreateOrUpdate(ctx context.Context, rule *domain.Rule) error {
	return m.CreateOrUpdateFunc(ctx, rule)
}

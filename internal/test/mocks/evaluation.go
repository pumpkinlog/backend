package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type EvaluationService struct {
	EvaluateRegionFunc  func(ctx context.Context, userID, regionID string) (*domain.RegionEvaluation, error)
	EvaluateRegionsFunc func(ctx context.Context, userID string) ([]*domain.RegionEvaluation, error)
}

func (m EvaluationService) EvaluateRegion(ctx context.Context, userID, regionID string) (*domain.RegionEvaluation, error) {
	return m.EvaluateRegionFunc(ctx, userID, regionID)
}

func (m EvaluationService) EvaluateRegions(ctx context.Context, userID string) ([]*domain.RegionEvaluation, error) {
	return m.EvaluateRegionsFunc(ctx, userID)
}

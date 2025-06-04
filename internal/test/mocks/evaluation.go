package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type EvaluationService struct {
	EvaluateRegionFunc  func(ctx context.Context, userID int64, regionID string) (*domain.RegionEvaluation, error)
	EvaluateRegionsFunc func(ctx context.Context, userID int64) ([]*domain.RegionEvaluation, error)
}

func (m EvaluationService) EvaluateRegion(ctx context.Context, userID int64, regionID string) (*domain.RegionEvaluation, error) {
	return m.EvaluateRegionFunc(ctx, userID, regionID)
}

func (m EvaluationService) EvaluateRegions(ctx context.Context, userID int64) ([]*domain.RegionEvaluation, error) {
	return m.EvaluateRegionsFunc(ctx, userID)
}

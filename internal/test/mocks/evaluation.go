package mocks

import (
	"context"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

type EvaluationService struct {
	EvaluationContextFunc func(ctx context.Context, userID int64, regionID domain.RegionID, pointInTime time.Time) (*domain.EvaluationContext, error)
	EvaluateRegionFunc    func(ctx context.Context, userID int64, regionID domain.RegionID, opts *domain.EvaluateOpts) (*domain.RegionEvaluation, error)
}

func (m EvaluationService) EvaluationContext(ctx context.Context, userID int64, regionID domain.RegionID, pointInTime time.Time) (*domain.EvaluationContext, error) {
	return m.EvaluationContextFunc(ctx, userID, regionID, pointInTime)
}

func (m EvaluationService) EvaluateRegion(ctx context.Context, userID int64, regionID domain.RegionID, opts *domain.EvaluateOpts) (*domain.RegionEvaluation, error) {
	return m.EvaluateRegionFunc(ctx, userID, regionID, opts)
}

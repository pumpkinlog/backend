package mocks

import (
	"github.com/pumpkinlog/backend/internal/app/engine"
)

type EvaluationService struct {
	EvaluateFunc func(params *engine.EvaluateRegionsParams) (*engine.RegionEvaluation, error)
}

func (m EvaluationService) Evaluate(params *engine.EvaluateRegionsParams) (*engine.RegionEvaluation, error) {
	return m.EvaluateFunc(params)
}

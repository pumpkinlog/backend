package strategies

import (
	"time"
)

type Strategy interface {
	Evaluate(config []byte, presences map[time.Time]struct{}) (StrategyEvaluation, error)
}

type StrategyEvaluation struct {
	Passed    bool           `json:"passed"`
	Count     int            `json:"count"`
	Remaining int            `json:"remaining"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

type Strategies struct {
	strategies map[string]Strategy
}

func NewStrategies() *Strategies {
	s := &Strategies{
		strategies: make(map[string]Strategy),
	}
	s.RegisterDefaultStrategies()
	return s
}

func (s *Strategies) Evaluate(rt string, cfg []byte, presences map[time.Time]struct{}) (StrategyEvaluation, error) {
	strategy, err := s.Strategy(rt)
	if err != nil {
		return StrategyEvaluation{}, err
	}

	evaluation, err := strategy.Evaluate(cfg, presences)
	if err != nil {
		return StrategyEvaluation{}, err
	}

	if evaluation.Metadata == nil {
		evaluation.Metadata = make(map[string]any)
	}

	return evaluation, nil
}

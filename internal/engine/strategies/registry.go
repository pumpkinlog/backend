package strategies

import (
	"fmt"
)

func (s *Strategies) RegisterDefaultStrategies() {
	s.RegisterStrategy("aggregate", &AggregateStrategy{})
	s.RegisterStrategy("average", &AverageStrategy{})
	s.RegisterStrategy("weighted", &WeightedStrategy{})
	s.RegisterStrategy("consecutive", &ConsecutiveStrategy{})
}

func (s *Strategies) Strategy(rt string) (Strategy, error) {
	strategy, exists := s.strategies[rt]
	if !exists {
		return nil, fmt.Errorf("no strategy registered for rule type %s", rt)
	}
	return strategy, nil
}

func (s *Strategies) RegisterStrategy(rt string, strat Strategy) {
	if _, exists := s.strategies[rt]; exists {
		panic(fmt.Sprintf("strategy for rule type %s already registered", rt))
	}
	s.strategies[rt] = strat
}

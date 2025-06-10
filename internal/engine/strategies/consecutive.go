package strategies

import "time"

// ConsecutiveStrategy counts consecutive presences within a fixed period.
type ConsecutiveStrategy struct {
	Threshold int `json:"threshold"`
}

func (s *ConsecutiveStrategy) Evaluate(config []byte, presences map[time.Time]struct{}) (StrategyEvaluation, error) {
	return StrategyEvaluation{}, nil
}

package strategies

import (
	"encoding/json"
	"fmt"
	"time"
)

// AggregateStrategy is a simple day-count threshold within a fixed period.
type AggregateStrategy struct {
	Threshold int `json:"threshold"`
}

func (s *AggregateStrategy) Evaluate(data []byte, presences map[time.Time]struct{}) (StrategyEvaluation, error) {
	if err := json.Unmarshal(data, &s); err != nil {
		return StrategyEvaluation{}, fmt.Errorf("invalid aggregate strategy config: %w", err)
	}

	count := len(presences)
	remaining := s.Threshold - count

	return StrategyEvaluation{
		Passed:    count > s.Threshold,
		Count:     count,
		Remaining: max(remaining, 0),
	}, nil
}

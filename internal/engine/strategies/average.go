package strategies

import (
	"encoding/json"
	"fmt"
	"time"
)

type AverageStrategy struct {
	Threshold int `json:"threshold"`
}

func (s *AverageStrategy) Evaluate(data []byte, presences map[time.Time]struct{}) (StrategyEvaluation, error) {
	if err := json.Unmarshal(data, &s); err != nil {
		return StrategyEvaluation{}, fmt.Errorf("invalid average strategy config: %w", err)
	}

	count := len(presences)
	ratio := float64(count) / float64(s.Threshold)

	return StrategyEvaluation{
		Passed:    count >= s.Threshold,
		Count:     count,
		Remaining: max(s.Threshold-count, 0),
		Metadata: map[string]any{
			"ratio": ratio,
		},
	}, nil
}

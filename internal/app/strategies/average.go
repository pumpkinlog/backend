package strategies

import (
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

type AverageStrategy struct {
	BaseStrategy
}

func NewAverageStrategy() RuleStrategy {
	return &AverageStrategy{}
}

func (s *AverageStrategy) Evaluate(rule *domain.Rule, presences []*domain.Presence) StrategyEvaluation {

	now := time.Now().UTC()
	start, end, _ := rule.Period(now)
	count := len(presences)
	avg := float64(count) / float64(rule.Threshold)
	remaining := rule.Threshold - count
	expiry := now.AddDate(0, 0, remaining)

	return StrategyEvaluation{
		Passed:         avg > float64(rule.Threshold),
		Count:          count,
		Remaining:      max(remaining, 0),
		Start:          start,
		End:            end,
		ConsecutiveEnd: time.Date(expiry.Year(), expiry.Month(), expiry.Day(), 23, 59, 59, 0, time.UTC),
		Metadata: map[string]any{
			"average": avg,
		},
	}
}

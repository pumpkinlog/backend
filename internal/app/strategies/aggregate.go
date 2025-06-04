package strategies

import (
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

type AggregateStrategy struct {
	BaseStrategy
}

func NewAggregateStrategy() RuleStrategy {
	return &AggregateStrategy{}
}

func (s *AggregateStrategy) Evaluate(rule *domain.Rule, presences []*domain.Presence) StrategyEvaluation {

	now := time.Now().UTC()
	start, end, _ := rule.Period(now)
	presences = s.filterPresencesByWindow(presences, start, end)
	count := len(presences)
	remaining := rule.Threshold - count
	expiry := now.AddDate(0, 0, remaining)

	return StrategyEvaluation{
		Passed:         count > rule.Threshold,
		Count:          count,
		Remaining:      max(remaining, 0),
		Start:          start,
		End:            end,
		ConsecutiveEnd: time.Date(expiry.Year(), expiry.Month(), expiry.Day(), 23, 59, 59, 0, time.UTC),
	}
}

package strategies

import (
	"time"

	"github.com/pumpkinlog/backend/internal/app/period"
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
	window := period.ComputeRulePeriod(now, rule)
	presences = s.filterPresencesByWindow(presences, window)
	count := len(presences)
	remaining := rule.Threshold - count
	expiry := now.AddDate(0, 0, remaining)

	return StrategyEvaluation{
		Passed:         count > rule.Threshold,
		Count:          count,
		Remaining:      max(remaining, 0),
		Start:          window.Start,
		End:            window.End,
		ConsecutiveEnd: time.Date(expiry.Year(), expiry.Month(), expiry.Day(), 23, 59, 59, 0, time.UTC),
	}
}

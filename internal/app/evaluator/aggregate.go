package evaluator

import (
	"time"

	"github.com/pumpkinlog/backend/internal/app/period"
	"github.com/pumpkinlog/backend/internal/domain"
)

func evaluateAggregateRule(rule *domain.Rule, presences []*domain.Presence) Evaluation {

	now := time.Now().UTC()
	window := period.ComputePeriod(now, rule)
	count := len(presences)
	remaining := rule.Threshold - count
	expiry := now.AddDate(0, 0, remaining)

	return Evaluation{
		Passed:         count > rule.Threshold,
		Count:          count,
		Remaining:      max(remaining, 0),
		Start:          window.Start,
		End:            window.End,
		ConsecutiveEnd: time.Date(expiry.Year(), expiry.Month(), expiry.Day(), 23, 59, 59, 0, time.UTC),
	}
}

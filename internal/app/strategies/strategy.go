package strategies

import (
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

// RuleStrategy defines the interface for different rule evaluation strategies
type RuleStrategy interface {
	// Evaluate evaluates a rule against a set of presences
	Evaluate(rule *domain.Rule, presences []*domain.Presence) StrategyEvaluation
}

type StrategyEvaluation struct {
	Passed         bool           `json:"passed"`
	Count          int            `json:"count"`
	Remaining      int            `json:"remaining"`
	Start          time.Time      `json:"start"`
	End            time.Time      `json:"end"`
	ConsecutiveEnd time.Time      `json:"consecutiveEnd"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

// BaseStrategy provides common functionality for all strategies
type BaseStrategy struct{}

// filterPresencesByWindow filters presences to only those within the given time window
func (s *BaseStrategy) filterPresencesByWindow(presences []*domain.Presence, start, end time.Time) []*domain.Presence {
	var filtered []*domain.Presence
	for _, p := range presences {
		if p.Date.After(start) && p.Date.Before(end) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

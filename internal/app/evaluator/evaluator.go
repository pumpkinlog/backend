package evaluator

import (
	"fmt"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

type Evaluation struct {
	Passed         bool           `json:"passed"`
	Count          int            `json:"count"`
	Remaining      int            `json:"remaining"`
	Start          time.Time      `json:"start"`
	End            time.Time      `json:"end"`
	ConsecutiveEnd time.Time      `json:"consecutiveEnd"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

type evaluateFunc = func(rule *domain.Rule, presences []*domain.Presence) Evaluation

var evaluators = map[domain.RuleType]evaluateFunc{
	domain.RuleTypeAggregate: evaluateAggregateRule,
	domain.RuleTypeAverage:   evaluateAverageRule,
	domain.RuleTypeWeighted:  evaluateWeightedRule,
}

func Evaluate(rule *domain.Rule, presences []*domain.Presence) (Evaluation, error) {

	if fn, ok := evaluators[rule.RuleType]; ok {
		return fn(rule, presences), nil
	}

	return Evaluation{}, fmt.Errorf("unsupported rule type: %s", rule.RuleType)
}

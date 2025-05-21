package engine

import (
	"fmt"

	"github.com/pumpkinlog/backend/internal/app/evaluator"
	"github.com/pumpkinlog/backend/internal/domain"
)

type RuleEvaluator interface {
	Evaluate(rule *domain.Rule, presences []*domain.Presence) evaluator.Evaluation
}

type RegionEvaluation struct {
	Region *domain.Region
	Rules  []RuleEvaluation
}

type RuleEvaluation struct {
	Passed     bool
	Rule       *domain.Rule
	Evaluation evaluator.Evaluation
	Conditions []ConditionEvaluation
}

type DefaultRuleEvaluator struct{}

func (e *DefaultRuleEvaluator) Evaluate(rule *domain.Rule, presences []*domain.Presence) evaluator.Evaluation {
	evaluation, err := evaluator.Evaluate(rule, presences)
	if err != nil {
		return evaluator.Evaluation{
			Metadata: map[string]any{
				"error": fmt.Errorf("skipped rule evaluation for %s: %w", rule.ID, err),
			},
		}
	}

	return evaluation
}

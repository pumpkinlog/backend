package engine

import (
	"fmt"

	"github.com/pumpkinlog/backend/internal/app/comparator"
	"github.com/pumpkinlog/backend/internal/app/strategies"
	"github.com/pumpkinlog/backend/internal/domain"
)

type Engine struct {
	strategies map[domain.RuleType]strategies.RuleStrategy
}

func NewEngine() *Engine {
	return &Engine{
		strategies: map[domain.RuleType]strategies.RuleStrategy{
			domain.RuleTypeAggregate: strategies.NewAggregateStrategy(),
			domain.RuleTypeAverage:   strategies.NewAverageStrategy(),
		},
	}
}

func (e *Engine) EvaluateRegion(agg *domain.RegionAggregate) ([]*domain.RuleEvaluation, bool, error) {
	evaluations := make([]*domain.RuleEvaluation, len(agg.Rules))
	var passed bool

	for i, rule := range agg.Rules {

		ruleConditions := make(map[int64]*domain.RuleCondition)
		for _, rc := range agg.RuleConditions {
			if rc.RuleID == rule.ID {
				ruleConditions[rc.ConditionID] = rc
			}
		}

		params := &EvaluateRuleParams{
			Presences:      agg.Presences,
			RuleConditions: ruleConditions,
			Conditions:     agg.Conditions,
			Answers:        agg.Answers,
		}

		evaluation, err := e.EvaluateRule(rule, params)
		if err != nil {
			return nil, false, fmt.Errorf("evaluate rule %d: %w", rule.ID, err)
		}

		passed = passed || evaluation.Passed

		evaluations[i] = evaluation
	}

	return evaluations, passed, nil
}

type EvaluateRuleParams struct {
	Presences      []*domain.Presence
	RuleConditions map[int64]*domain.RuleCondition
	Conditions     map[int64]*domain.Condition
	Answers        map[int64]*domain.Answer
}

func (e *Engine) EvaluateRule(rule *domain.Rule, params *EvaluateRuleParams) (*domain.RuleEvaluation, error) {
	evaluation := &domain.RuleEvaluation{
		RuleID: rule.ID,
		Passed: true,
	}

	for conditionID, rc := range params.RuleConditions {
		condition, ok := params.Conditions[conditionID]
		if !ok {
			return nil, fmt.Errorf("condition %d not found for rule %d", conditionID, rule.ID)
		}
		answer := params.Answers[conditionID]

		ce, err := e.EvaluateCondition(condition, rc, answer)
		if err != nil {
			return nil, fmt.Errorf("evaluate condition %d: %w", conditionID, err)
		}

		evaluation.ConditionEvaluations = append(evaluation.ConditionEvaluations, ce)

		if !ce.Passed && evaluation.Passed {
			evaluation.Passed = false
		}
	}

	if !evaluation.Passed {
		return evaluation, nil
	}

	se, err := e.evaluateStrategy(rule, params.Presences)
	if err != nil {
		return nil, fmt.Errorf("evaluate rule strategy: %w", err)
	}

	evaluation.Passed = se.Passed
	evaluation.Count = &se.Count
	evaluation.Remaining = &se.Remaining
	evaluation.Start = se.Start
	evaluation.End = se.End
	evaluation.ConsecutiveEnd = se.ConsecutiveEnd
	evaluation.Metadata = se.Metadata

	return evaluation, nil
}

func (e *Engine) EvaluateCondition(cond *domain.Condition, rc *domain.RuleCondition, ans *domain.Answer) (*domain.ConditionEvaluation, error) {
	evaluation := &domain.ConditionEvaluation{
		ConditionID: cond.ID,
	}

	if ans == nil {
		evaluation.Skipped = true
		return evaluation, nil
	}

	passed, err := comparator.Compare(cond.Type, rc.Comparator, rc.Expected, ans.Value)
	if err != nil {
		return nil, fmt.Errorf("compare values: %w", err)
	}

	evaluation.Passed = passed
	evaluation.Expected = rc.Expected
	evaluation.Actual = ans.Value
	evaluation.Comparator = rc.Comparator

	return evaluation, nil
}

func (e *Engine) evaluateStrategy(rule *domain.Rule, presences []*domain.Presence) (strategies.StrategyEvaluation, error) {
	strategy, ok := e.strategies[rule.RuleType]
	if !ok {
		return strategies.StrategyEvaluation{}, fmt.Errorf("unsupported rule type: %s", rule.RuleType)
	}

	return strategy.Evaluate(rule, presences), nil
}

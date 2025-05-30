package engine

import (
	"fmt"
	"time"

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

type EvaluateRegionParams struct {
	Region         *domain.Region
	Presences      []*domain.Presence
	Rules          []*domain.Rule
	Conditions     map[string]*domain.Condition       // conditionID -> condition
	Answers        map[string]*domain.Answer          // conditionID -> answer
	RuleConditions map[string][]*domain.RuleCondition // ruleID -> ruleConditions
}

func (e *Engine) EvaluateRegion(params *EvaluateRegionParams) (*domain.RegionEvaluation, error) {
	evaluation := &domain.RegionEvaluation{
		Region:          params.Region,
		RuleEvaluations: make([]domain.RuleEvaluation, len(params.Rules)),
		EvaluatedAt:     time.Now().UTC(),
		Conditions:      params.Conditions,
		Answers:         params.Answers,
	}

	for i, rule := range params.Rules {
		re, err := e.EvaluateRule(rule, params)
		if err != nil {
			return nil, fmt.Errorf("evaluate rule %s: %w", rule.ID, err)
		}

		evaluation.RuleEvaluations[i] = re
	}

	return evaluation, nil
}

func (e *Engine) EvaluateRule(rule *domain.Rule, params *EvaluateRegionParams) (domain.RuleEvaluation, error) {
	ruleConditions := params.RuleConditions[rule.ID]
	evaluation := domain.RuleEvaluation{
		Passed:               true,
		Rule:                 rule,
		ConditionEvaluations: make([]domain.ConditionEvaluation, len(ruleConditions)),
	}

	for i, rc := range ruleConditions {
		cond := params.Conditions[rc.ConditionID]
		answer := params.Answers[cond.ID]

		ce, err := e.evaluateCondition(cond, rc, answer)
		if err != nil {
			return domain.RuleEvaluation{}, fmt.Errorf("evaluate condition %s: %w", cond.ID, err)
		}

		evaluation.ConditionEvaluations[i] = ce

		if !ce.Passed {
			evaluation.Passed = false
		}
	}

	if !evaluation.Passed {
		return evaluation, nil
	}

	se, err := e.evaluateRule(rule, params.Presences)
	if err != nil {
		return domain.RuleEvaluation{}, fmt.Errorf("evaluate rule strategy: %w", err)
	}

	evaluation.Passed = se.Passed
	evaluation.Count = se.Count
	evaluation.Remaining = se.Remaining
	evaluation.Start = se.Start
	evaluation.End = se.End
	evaluation.ConsecutiveEnd = se.ConsecutiveEnd
	evaluation.Metadata = se.Metadata

	return evaluation, nil
}

func (e *Engine) evaluateCondition(cond *domain.Condition, rc *domain.RuleCondition, ans *domain.Answer) (domain.ConditionEvaluation, error) {
	evaluation := domain.ConditionEvaluation{
		ConditionID: cond.ID,
	}

	if ans == nil {
		evaluation.Skipped = true
		return evaluation, nil
	}

	valueType := comparator.Type(cond.Type)
	operator := comparator.Operator(rc.Comparator)

	passed, err := comparator.Compare(valueType, operator, rc.Expected, ans.Value)
	if err != nil {
		return domain.ConditionEvaluation{}, fmt.Errorf("compare values: %w", err)
	}

	evaluation.Passed = passed
	evaluation.Expected = rc.Expected
	evaluation.Actual = ans.Value

	return evaluation, nil
}

func (e *Engine) evaluateRule(rule *domain.Rule, presences []*domain.Presence) (strategies.StrategyEvaluation, error) {
	strategy, ok := e.strategies[rule.RuleType]
	if !ok {
		return strategies.StrategyEvaluation{}, fmt.Errorf("unsupported rule type: %s", rule.RuleType)
	}

	return strategy.Evaluate(rule, presences), nil
}

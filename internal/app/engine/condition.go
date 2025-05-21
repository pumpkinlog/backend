package engine

import "github.com/pumpkinlog/backend/internal/domain"

type ConditionEvaluator interface {
	Evaluate(condition *domain.Condition, answer *domain.Answer) ConditionEvaluation
}

type ConditionEvaluation struct {
	Passed    bool
	Skipped   bool
	Condition *domain.Condition
	Answer    *domain.Answer
}

type DefaultConditionEvaluator struct{}

func (e *DefaultConditionEvaluator) Evaluate(condition *domain.Condition, answer *domain.Answer) ConditionEvaluation {
	evaluation := ConditionEvaluation{
		Condition: condition,
		Answer:    answer,
	}

	if answer == nil {
		evaluation.Skipped = true
		return evaluation
	}

	evaluation.Passed = e.evaluateLogic(condition.Type, condition.Comparator, condition.Expected, answer.Value)

	return evaluation
}

func (e *DefaultConditionEvaluator) evaluateLogic(condType domain.ConditionType, cmp domain.Comparator, expected, actual any) bool {
	return compare(cmp, expected, actual)
}

func compare(cmp domain.Comparator, expected, actual any) bool {
	switch e := expected.(type) {
	case int:
		if a, ok := actual.(int); ok {
			return compareOrdered(e, a, cmp)
		}
	case float32:
		if a, ok := actual.(float32); ok {
			return compareOrdered(e, a, cmp)
		}
	case string:
		if a, ok := actual.(string); ok {
			return compareOrdered(e, a, cmp)
		}
	}
	return false
}

func compareOrdered[T string | int | float32](expected, actual T, cmp domain.Comparator) bool {
	switch cmp {
	case domain.ComparatorEquals:
		return expected == actual
	case domain.ComparatorNotEquals:
		return expected != actual
	case domain.ComparatorGreaterThan:
		return expected > actual
	case domain.ComparatorLessThan:
		return expected < actual
	case domain.ComparatorGreaterThanOrEquals:
		return expected >= actual
	case domain.ComparatorLessThanOrEquals:
		return expected <= actual
	default:
		return false
	}
}

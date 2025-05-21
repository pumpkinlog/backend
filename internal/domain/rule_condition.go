package domain

import "context"

type RuleCondition struct {
	RuleID      string
	ConditionID string
}

type RuleConditionRepository interface {
	Link(ctx context.Context, ruleID, conditionID string) error
	Unlink(ctx context.Context, ruleID, conditionID string) error
}

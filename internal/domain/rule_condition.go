package domain

import "context"

type Comparator string

const (
	ComparatorEquals              Comparator = "eq"
	ComparatorNotEquals           Comparator = "neq"
	ComparatorGreaterThan         Comparator = "gt"
	ComparatorGreaterThanOrEquals Comparator = "gte"
	ComparatorLessThan            Comparator = "lt"
	ComparatorLessThanOrEquals    Comparator = "lte"
	ComparatorContains            Comparator = "contains"
	ComparatorIn                  Comparator = "in"
	ComparatorNotIn               Comparator = "not_in"
)

type RuleCondition struct {
	RuleID      string
	ConditionID string
	Weight      int
	Comparator  Comparator
	Expected    any
}

type RuleConditionFilter struct {
	RuleIDs []string
}

type RuleConditionRepository interface {
	GetByID(ctx context.Context, ruleID, conditionID string) (*RuleCondition, error)
	List(ctx context.Context, filter *RuleConditionFilter) ([]*RuleCondition, error)
	CreateOrUpdate(ctx context.Context, rc *RuleCondition) error
	Delete(ctx context.Context, ruleID, conditionID string) error
}

package domain

import (
	"context"
)

type (
	ConditionType string
	Comparator    string
)

const (
	ConditionTypeBoolean     ConditionType = "boolean"
	ConditionTypeInteger     ConditionType = "integer"
	ConditionTypeSelect      ConditionType = "select"
	ConditionTypeMultiSelect ConditionType = "multi_select"

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

type Condition struct {
	ID         string        `json:"id"`
	RuleID     string        `json:"ruleId"`
	Prompt     string        `json:"prompt"`
	Type       ConditionType `json:"type"`
	Comparator Comparator    `json:"-"`
	Expected   any           `json:"-"`
}

type ConditionFilter struct {
	RuleIDs []string
	Page    *int
	Limit   *int
}

type ConditionService interface {
	Create(ctx context.Context, condition *Condition, ruleIDs []string) error
	Delete(ctx context.Context, conditionID, ruleID string) error
}

type ConditionRepository interface {
	GetByID(ctx context.Context, conditionID string) (*Condition, error)
	List(ctx context.Context, filter *ConditionFilter) ([]*Condition, error)

	CreateOrUpdate(ctx context.Context, condition *Condition) error

	Delete(ctx context.Context, conditionID string) error
}

package domain

import (
	"context"
)

type ConditionType string

const (
	ConditionTypeString      ConditionType = "string"
	ConditionTypeBoolean     ConditionType = "boolean"
	ConditionTypeInteger     ConditionType = "integer"
	ConditionTypeSelect      ConditionType = "select"
	ConditionTypeMultiSelect ConditionType = "multi_select"
)

type Condition struct {
	ID     string        `json:"id"`
	Prompt string        `json:"prompt"`
	Type   ConditionType `json:"type"`
}

type ConditionService interface {
	Create(ctx context.Context, condition *Condition, ruleIDs []string) error
	Delete(ctx context.Context, conditionID, ruleID string) error
}

type ConditionFilter struct {
	ConditionIDs []string
	Page         *int
	Limit        *int
}

type ConditionRepository interface {
	GetByID(ctx context.Context, conditionID string) (*Condition, error)
	List(ctx context.Context, filter *ConditionFilter) ([]*Condition, error)
	CreateOrUpdate(ctx context.Context, condition *Condition) error
	Delete(ctx context.Context, conditionID string) error
}

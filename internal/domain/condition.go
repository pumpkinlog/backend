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
	ID       Code          `json:"id"`
	RegionID RegionID      `json:"regionId"`
	Prompt   string        `json:"prompt"`
	Type     ConditionType `json:"type"`
}

func (t ConditionType) Valid() bool {
	switch t {
	case ConditionTypeString, ConditionTypeBoolean, ConditionTypeInteger, ConditionTypeSelect, ConditionTypeMultiSelect:
		return true
	default:
		return false
	}
}

func (c *Condition) Validate() error {
	if err := c.ID.Validate(); err != nil {
		return err
	}

	if err := c.RegionID.Validate(); err != nil {
		return err
	}

	if c.Prompt == "" {
		return ValidationError("prompt is required")
	}

	if !c.Type.Valid() {
		return ValidationError("invalid condition type: %s", c.Type)
	}

	return nil
}

type ConditionService interface {
	GetByID(ctx context.Context, conditionID Code) (*Condition, error)
	List(ctx context.Context, filter *ConditionFilter) ([]*Condition, error)
	CreateOrUpdate(ctx context.Context, condition *Condition) error
}

type ConditionFilter struct {
	RegionIDs []RegionID
}

type ConditionRepository interface {
	GetByID(ctx context.Context, conditionID Code) (*Condition, error)
	List(ctx context.Context, filter *ConditionFilter) ([]*Condition, error)
	ListByRegionID(ctx context.Context, regionID RegionID) ([]*Condition, error)
	CreateOrUpdate(ctx context.Context, condition *Condition) error
	Delete(ctx context.Context, conditionID Code) error
}

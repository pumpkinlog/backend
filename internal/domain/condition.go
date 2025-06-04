package domain

import (
	"context"
	"fmt"
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
	ID       int64         `json:"id"`
	RegionID string        `json:"regionId"`
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

	if c.RegionID == "" {
		return fmt.Errorf("regionId cannot be empty")
	}

	if c.Prompt == "" {
		return fmt.Errorf("prompt cannot be empty")
	}

	if !c.Type.Valid() {
		return fmt.Errorf("invalid condition type: %s", c.Type)
	}

	return nil
}

type ConditionService interface {
	CreateOrUpdate(ctx context.Context, condition *Condition) error
}

type ConditionFilter struct {
	RegionIDs []string
	Page      *int
	Limit     *int
}

type ConditionRepository interface {
	GetByID(ctx context.Context, conditionID int64) (*Condition, error)
	List(ctx context.Context, filter *ConditionFilter) ([]*Condition, error)
	ListByRegionID(ctx context.Context, regionID string) ([]*Condition, error)
	CreateOrUpdate(ctx context.Context, condition *Condition) error
	Delete(ctx context.Context, conditionID int64) error
}

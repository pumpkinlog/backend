package domain

import (
	"context"
	"fmt"
)

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
	RuleID      int64
	ConditionID int64
	RegionID    string
	Comparator  Comparator
	Expected    any
}

func (rc *RuleCondition) Validate() error {

	if rc.RegionID == "" {
		return fmt.Errorf("regionId cannot be empty")
	}

	if rc.Comparator == "" {
		return fmt.Errorf("comparator cannot be empty")
	}

	if rc.Expected == nil {
		return fmt.Errorf("expected value cannot be nil")
	}

	if rc.RuleID <= 0 {
		return fmt.Errorf("ruleId must be greater than 0")
	}

	if rc.ConditionID <= 0 {
		return fmt.Errorf("conditionId must be greater than 0")
	}

	return nil
}

type RuleConditionService interface {
	CreateOrUpdate(ctx context.Context, rc *RuleCondition) error
}

type RuleConditionRepository interface {
	GetByID(ctx context.Context, ruleID, conditionID int64) (*RuleCondition, error)
	ListByRegionID(ctx context.Context, regionID string) ([]*RuleCondition, error)
	CreateOrUpdate(ctx context.Context, rc *RuleCondition) error
	Delete(ctx context.Context, ruleID, conditionID int64) error
}

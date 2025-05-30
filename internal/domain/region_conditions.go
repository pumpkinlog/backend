package domain

import "context"

type RegionCondition struct {
	RegionID    string `json:"regionId"`
	ConditionID string `json:"conditionId"`
}

type RegionConditionFilter struct {
	Page  int
	Limit int
}

type RegionConditionRepository interface {
	GetByID(ctx context.Context, regionID, conditionID string) (*RegionCondition, error)
	List(ctx context.Context, filter *RegionConditionFilter) ([]*RegionCondition, error)
	CreateOrUpdate(ctx context.Context, rc *RegionCondition) error
	Delete(ctx context.Context, regionID, conditionID string) error
}

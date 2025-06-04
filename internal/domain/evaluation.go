package domain

import (
	"context"
	"time"
)

// RegionEvaluation represents the result of evaluating a region against its rules
type RegionEvaluation struct {
	UserID          int64             `json:"userId"`
	RegionID        string            `json:"regionId"`
	Passed          bool              `json:"passed"`
	RuleEvaluations []*RuleEvaluation `json:"ruleEvaluations"`
	EvaluatedAt     time.Time         `json:"evaluatedAt"`
}

// RuleEvaluation represents the result of evaluating a single rule
type RuleEvaluation struct {
	RuleID               int64                  `json:"ruleId"`
	Passed               bool                   `json:"passed"`
	Count                *int                   `json:"count,omitempty"`
	Remaining            *int                   `json:"remaining,omitempty"`
	Start                time.Time              `json:"start,omitzero"`
	End                  time.Time              `json:"end,omitzero"`
	ConsecutiveEnd       time.Time              `json:"consecutiveEnd,omitzero"`
	Metadata             map[string]any         `json:"metadata,omitempty"`
	ConditionEvaluations []*ConditionEvaluation `json:"conditionEvaluations,omitempty"`
}

// ConditionEvaluation represents the result of evaluating a single condition
type ConditionEvaluation struct {
	ConditionID int64      `json:"conditionId"`
	Passed      bool       `json:"passed"`
	Skipped     bool       `json:"skipped"`
	Comparator  Comparator `json:"comparator,omitempty"`
	Expected    any        `json:"expected,omitempty"`
	Actual      any        `json:"actual,omitempty"`
}

type EvaluationService interface {
	EvaluateRegion(ctx context.Context, userID int64, regionID string) (*RegionEvaluation, error)
}

type EvaluationRepository interface {
	GetByID(ctx context.Context, userID int64, regionID string) (*RegionEvaluation, error)
	List(ctx context.Context, userID int64) ([]*RegionEvaluation, error)
	CreateOrUpdate(ctx context.Context, evaluation *RegionEvaluation) error
	DeleteByRegionID(ctx context.Context, regionID string) error
}

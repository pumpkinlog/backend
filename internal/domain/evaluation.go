package domain

import (
	"context"
	"time"
)

type ConditionEvaluation struct {
	Condition *Condition `json:"condition"`
	Answer    *Answer    `json:"answer,omitzero"`
	Passed    bool       `json:"passed"`
	Skipped   bool       `json:"skipped"`
	Expected  any        `json:"expected,omitempty"`
	Actual    any        `json:"actual,omitempty"`
	Reason    string     `json:"reason,omitzero"`
}

type RuleLogicEvaluation struct {
	Resident       bool           `json:"resident"`
	Count          int            `json:"count"`
	Remaining      int            `json:"remaining"`
	Start          time.Time      `json:"start"`
	End            time.Time      `json:"end"`
	ConsecutiveEnd time.Time      `json:"consecutiveEnd"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

type RuleEvaluation struct {
	Passed     bool                  `json:"passed"`
	Rule       *Rule                 `json:"rule"`
	Logic      RuleLogicEvaluation   `json:"ruleEvaluation,omitzero"`
	Conditions []ConditionEvaluation `json:"conditionEvaluations,omitempty"`
}

type RegionEvaluation struct {
	UserID      string           `json:"userId"`
	RegionID    string           `json:"regionId"`
	Passed      bool             `json:"passed"`
	Evaluations []RuleEvaluation `json:"evaluations"`
	EvaluatedAt time.Time        `json:"evaluatedAt"`
}

type EvaluationService interface {
	EvaluateRegion(ctx context.Context, userID, regionID string) (*RegionEvaluation, error)
	EvaluateRegions(ctx context.Context, userID string, regionIDs []string) ([]*RegionEvaluation, error)
	EvaluateAllRegions(ctx context.Context, userID string) ([]*RegionEvaluation, error)
}

type EvaluationRepository interface {
	GetByID(ctx context.Context, userID, regionID string) (*RegionEvaluation, error)
	List(ctx context.Context, userID string) ([]*RegionEvaluation, error)

	CreateOrUpdate(ctx context.Context, evaluation *RegionEvaluation) error
}

package domain

import (
	"context"
	"time"
)

// RegionEvaluation represents the result of evaluating a region against its rules
type RegionEvaluation struct {
	UserID          string                `json:"userId"`
	RegionID        string                `json:"regionId"`
	Passed          bool                  `json:"passed"`
	Region          *Region               `json:"region"`
	EvaluatedAt     time.Time             `json:"evaluatedAt"`
	RuleEvaluations []RuleEvaluation      `json:"ruleEvaluations"`
	Conditions      map[string]*Condition `json:"conditions"`
	Answers         map[string]*Answer    `json:"answers"`
}

// RuleEvaluation represents the result of evaluating a single rule
type RuleEvaluation struct {
	Passed               bool                  `json:"passed"`
	Count                int                   `json:"count"`
	Remaining            int                   `json:"remaining"`
	Start                time.Time             `json:"start"`
	End                  time.Time             `json:"end"`
	ConsecutiveEnd       time.Time             `json:"consecutiveEnd"`
	Metadata             map[string]any        `json:"metadata,omitempty"`
	Rule                 *Rule                 `json:"rule"`
	ConditionEvaluations []ConditionEvaluation `json:"conditionEvaluations"`
}

// ConditionEvaluation represents the result of evaluating a single condition
type ConditionEvaluation struct {
	ConditionID string `json:"conditionId"`
	Passed      bool   `json:"passed"`
	Skipped     bool   `json:"skipped"`
	Expected    any    `json:"expected,omitempty"`
	Actual      any    `json:"actual,omitempty"`
}

type EvaluationService interface {
	EvaluateRegion(ctx context.Context, userID, regionID string) (*RegionEvaluation, error)
	EvaluateRegions(ctx context.Context, userID string) ([]*RegionEvaluation, error)
}

type EvaluationRepository interface {
	GetByID(ctx context.Context, userID, regionID string) (*RegionEvaluation, error)
	List(ctx context.Context, userID string) ([]*RegionEvaluation, error)
	CreateOrUpdate(ctx context.Context, evaluation *RegionEvaluation) error
}

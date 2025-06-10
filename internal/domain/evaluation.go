package domain

import (
	"context"
	"time"
)

type EvaluationContext struct {
	At        time.Time
	Region    *Region
	Presences []*Presence
	Rules     []*Rule
	Answers   map[Code]*Answer
}

type RegionEvaluation struct {
	RegionID    RegionID              `json:"regionId"`
	UserID      int64                 `json:"userId"`
	Passed      bool                  `json:"passed"`
	Nodes       []EvaluationComponent `json:"nodes"`
	PointInTime time.Time             `json:"pointInTime"`
	EvaluatedAt time.Time             `json:"evaluatedAt"`
}

type (
	EvaluationStatus string
	ComponentType    string
)

const (
	EvaluationStatusEvaluated  EvaluationStatus = "evaluated"
	EvaluationStatusUnanswered EvaluationStatus = "unanswered"
	EvaluationStatusError      EvaluationStatus = "error"

	ComponentTypeComposite ComponentType = "composite"
	ComponentTypeStrategy  ComponentType = "strategy"
	ComponentTypeCondition ComponentType = "condition"
)

type EvaluationComponent interface {
	IsPassed() bool
}

type RuleEvaluation struct {
	EvaluationComponent
	RuleID string `json:"ruleId"`
}

type CompositeEvaluation struct {
	NodeType   NodeType              `json:"nodeType"`
	Status     EvaluationStatus      `json:"status"`
	Passed     bool                  `json:"passed"`
	Components []EvaluationComponent `json:"components"`
}

func (e *CompositeEvaluation) IsPassed() bool {
	return e.Passed
}

type StrategyEvaluation struct {
	Type      ComponentType    `json:"type"`
	Strategy  string           `json:"strategy"`
	Passed    bool             `json:"passed"`
	Status    EvaluationStatus `json:"status"`
	Reason    string           `json:"reason,omitempty"`
	Start     time.Time        `json:"start"`
	End       time.Time        `json:"end"`
	Count     int              `json:"count"`
	Remaining int              `json:"remaining"`
}

func (e *StrategyEvaluation) IsPassed() bool {
	return e.Passed
}

type ConditionEvaluation struct {
	Type        ComponentType    `json:"type"`
	ConditionID Code             `json:"conditionId"`
	Expected    any              `json:"expected"`
	Actual      any              `json:"actual"`
	Comparator  Comparator       `json:"comparator"`
	Status      EvaluationStatus `json:"status"`
	Passed      bool             `json:"passed"`
	Reason      string           `json:"reason,omitempty"`
}

func (e *ConditionEvaluation) IsPassed() bool {
	return e.Passed
}

type EvaluateOpts struct {
	// PointInTime is the time at which to evaluate the region.
	PointInTime time.Time
	// ForceRecompute indicates whether to force recompute the evaluation even if it already exists.
	Recompute bool
	// Cache indicates whether to cache the evaluation result.
	Cache bool
	// Publish indicates whether to publish the evaluation result to message queues.
	Publish bool
}

type EvaluationService interface {
	EvaluateRegion(ctx context.Context, userID int64, regionID RegionID, opts *EvaluateOpts) (*RegionEvaluation, error)
}

type EvaluationRepository interface {
	GetByID(ctx context.Context, userID int64, regionID RegionID) (*RegionEvaluation, error)
	List(ctx context.Context, userID int64) ([]*RegionEvaluation, error)
	CreateOrUpdate(ctx context.Context, evaluation *RegionEvaluation) error
	DeleteByUserAndRegionID(ctx context.Context, userID int64, regionID RegionID) error
	DeleteByRegionID(ctx context.Context, regionID RegionID) error
}

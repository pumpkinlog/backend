package domain

import (
	"context"
	"encoding/json"
)

type NodeType string

const (
	NodeTypeCompositeAnd NodeType = "and"
	NodeTypeCompositeAny NodeType = "any"
	NodeTypeStrategy     NodeType = "strategy"
	NodeTypeCondition    NodeType = "condition"
)

type RuleNode struct {
	Type  NodeType        `json:"type"`
	Props json.RawMessage `json:"props"`
}

func (rn *RuleNode) Validate() error {
	if rn.Type == "" {
		return ValidationError("type is required")
	}

	if rn.Props == nil {
		return ValidationError("props is required")
	}

	return nil
}

func (rn *RuleNode) Unmarshal(dst any) error {
	return json.Unmarshal(rn.Props, dst)
}

type Rule struct {
	ID          Code     `json:"id"`
	RegionID    RegionID `json:"regionId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Node        RuleNode `json:"node"`
}

func (r *Rule) Validate() error {
	if err := r.ID.Validate(); err != nil {
		return err
	}

	if err := r.RegionID.Validate(); err != nil {
		return err
	}

	if r.Name == "" {
		return ValidationError("name is required")
	}

	if r.Description == "" {
		return ValidationError("description is required")
	}

	if err := r.Node.Validate(); err != nil {
		return err
	}

	return nil
}

type Operator string

const (
	RuleOperatorAnd Operator = "and"
	RuleOperatorOr  Operator = "or"
)

type CompositeNode struct {
	Operator Operator   `json:"operator"`
	Nodes    []RuleNode `json:"nodes"`
}

func (n *CompositeNode) Validate() error {
	switch n.Operator {
	case RuleOperatorAnd, RuleOperatorOr:
	default:
		return ValidationError("composite node has invalid operator: %s", n.Operator)
	}

	if len(n.Nodes) < 2 {
		return ValidationError("operator %s requires at least 2 child nodes", n.Operator)
	}

	for _, cn := range n.Nodes {
		if err := cn.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type EvaluatorNode struct {
	Type   string          `json:"type"`
	Period Period          `json:"period"`
	Props  json.RawMessage `json:"props,omitempty"`
}

func (n *EvaluatorNode) Validate() error {
	if n.Type == "" {
		return ValidationError("strategy node type cannot be empty")
	}

	if err := n.Period.Validate(); err != nil {
		return err
	}

	return nil
}

type Comparator string

const (
	ComparatorEquals    Comparator = "eq"
	ComparatorNotEquals Comparator = "neq"
)

type ConditionNode struct {
	ConditionID Code       `json:"conditionId"`
	Equals      any        `json:"equals"`
	Comparator  Comparator `json:"comparator"`
}

func (n *ConditionNode) Validate() error {
	if err := n.ConditionID.Validate(); err != nil {
		return err
	}

	switch n.Comparator {
	case ComparatorEquals, ComparatorNotEquals:
	default:
		return ValidationError("unsupported comparator: %s", n.Comparator)
	}

	return nil
}

type PeriodType string

const (
	PeriodTypeYear    PeriodType = "year"
	PeriodTypeRolling PeriodType = "rolling"
)

type Period struct {
	Type          PeriodType `json:"type"`
	OffsetYears   int        `json:"offsetYears"`
	Years         int        `json:"years"`
	RollingDays   int        `json:"rollingDays"`
	RollingMonths int        `json:"rollingMonths"`
	RollingYears  int        `json:"rollingYears"`
}

func (p *Period) Validate() error {
	switch p.Type {
	case PeriodTypeYear:
		if p.Years <= 0 {
			return ValidationError("year period must have years greater than 0")
		}
	case PeriodTypeRolling:
		if p.RollingDays <= 0 && p.RollingMonths <= 0 && p.RollingYears <= 0 {
			return ValidationError("rolling period must have one days/months/years greater than 0")
		}
	default:
		return ValidationError("unknown period type: %s", p.Type)
	}

	return nil
}

type RuleService interface {
	GetByID(ctx context.Context, ruleID Code) (*Rule, error)
	List(ctx context.Context, filter *RuleFilter) ([]*Rule, error)
	CreateOrUpdate(ctx context.Context, rule *Rule) error
}

type RuleFilter struct {
	RegionIDs []RegionID
}

type RuleRepository interface {
	GetByID(ctx context.Context, ruleID Code) (*Rule, error)
	List(ctx context.Context, filter *RuleFilter) ([]*Rule, error)
	ListByRegionID(ctx context.Context, regionID RegionID) ([]*Rule, error)
	CreateOrUpdate(ctx context.Context, rule *Rule) error
}

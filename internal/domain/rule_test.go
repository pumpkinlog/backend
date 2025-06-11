package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRuleNode(t *testing.T) {
	baseProps := json.RawMessage(`{"foo":"bar"}`)
	baseNode := RuleNode{
		Type:  NodeTypeCompositeAnd,
		Props: baseProps,
	}

	tests := []struct {
		name    string
		modify  func(rn RuleNode) RuleNode
		wantErr error
	}{
		{
			name:   "valid RuleNode",
			modify: func(rn RuleNode) RuleNode { return rn },
		},
		{
			name: "missing type",
			modify: func(rn RuleNode) RuleNode {
				rn.Type = ""
				return rn
			},
			wantErr: ValidationError("type is required"),
		},
		{
			name: "missing props",
			modify: func(rn RuleNode) RuleNode {
				rn.Props = nil
				return rn
			},
			wantErr: ValidationError("props is required"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rn := tc.modify(baseNode)
			err := rn.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRule(t *testing.T) {
	baseRule := Rule{
		ID:          Code("VALID_ID"),
		RegionID:    RegionID("GB"),
		Name:        "Test Rule",
		Description: "A test rule",
		Node: RuleNode{
			Type:  NodeTypeCompositeAnd,
			Props: json.RawMessage(`{"foo":"bar"}`),
		},
	}

	tests := []struct {
		name    string
		modify  func(r Rule) Rule
		wantErr error
	}{
		{
			name:   "valid rule",
			modify: func(r Rule) Rule { return r },
		},
		{
			name: "invalid ID",
			modify: func(r Rule) Rule {
				r.ID = Code("")
				return r
			},
			wantErr: ValidationError("code is required"),
		},
		{
			name: "invalid RegionID",
			modify: func(r Rule) Rule {
				r.RegionID = RegionID("")
				return r
			},
			wantErr: ValidationError("region ID is required"),
		},
		{
			name: "missing name",
			modify: func(r Rule) Rule {
				r.Name = ""
				return r
			},
			wantErr: ValidationError("name is required"),
		},
		{
			name: "missing description",
			modify: func(r Rule) Rule {
				r.Description = ""
				return r
			},
			wantErr: ValidationError("description is required"),
		},
		{
			name: "invalid node",
			modify: func(r Rule) Rule {
				r.Node.Type = ""
				return r
			},
			wantErr: ValidationError("type is required"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.modify(baseRule)
			err := r.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateCompositeNode(t *testing.T) {
	validNode := RuleNode{
		Type:  NodeTypeCondition,
		Props: json.RawMessage(`{"conditionId":"valid"}`),
	}

	baseComposite := CompositeNode{
		Operator: RuleOperatorAnd,
		Nodes:    []RuleNode{validNode, validNode},
	}

	tests := []struct {
		name    string
		modify  func(cn CompositeNode) CompositeNode
		wantErr error
	}{
		{
			name:   "valid composite node",
			modify: func(cn CompositeNode) CompositeNode { return cn },
		},
		{
			name: "invalid operator",
			modify: func(cn CompositeNode) CompositeNode {
				cn.Operator = "invalid"
				return cn
			},
			wantErr: ValidationError("composite node has invalid operator: invalid"),
		},
		{
			name: "less than two child nodes",
			modify: func(cn CompositeNode) CompositeNode {
				cn.Nodes = []RuleNode{validNode}
				return cn
			},
			wantErr: ValidationError("operator and requires at least 2 child nodes"),
		},
		{
			name: "invalid child node",
			modify: func(cn CompositeNode) CompositeNode {
				badNode := RuleNode{Type: "", Props: nil}
				cn.Nodes = []RuleNode{validNode, badNode}
				return cn
			},
			wantErr: ValidationError("type is required"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cn := tc.modify(baseComposite)
			err := cn.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateEvaluatorNode(t *testing.T) {
	validPeriod := Period{
		Type:  PeriodTypeYear,
		Years: 1,
	}

	baseEvaluator := EvaluatorNode{
		Type:   "some-strategy",
		Period: validPeriod,
		Props:  json.RawMessage(`{"foo":"bar"}`),
	}

	tests := []struct {
		name    string
		modify  func(en EvaluatorNode) EvaluatorNode
		wantErr error
	}{
		{
			name:   "valid evaluator node",
			modify: func(en EvaluatorNode) EvaluatorNode { return en },
		},
		{
			name: "empty type",
			modify: func(en EvaluatorNode) EvaluatorNode {
				en.Type = ""
				return en
			},
			wantErr: ValidationError("strategy node type cannot be empty"),
		},
		{
			name: "invalid period",
			modify: func(en EvaluatorNode) EvaluatorNode {
				en.Period.Type = "invalid"
				return en
			},
			wantErr: ValidationError("unknown period type: invalid"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			en := tc.modify(baseEvaluator)
			err := en.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateConditionNode(t *testing.T) {
	baseCond := ConditionNode{
		ConditionID: Code("VALID_CODE"),
		Comparator:  ComparatorEquals,
		Equals:      "some-value",
	}

	tests := []struct {
		name    string
		modify  func(cn ConditionNode) ConditionNode
		wantErr error
	}{
		{
			name:   "valid condition node",
			modify: func(cn ConditionNode) ConditionNode { return cn },
		},
		{
			name: "invalid ConditionID",
			modify: func(cn ConditionNode) ConditionNode {
				cn.ConditionID = Code("")
				return cn
			},
			wantErr: ValidationError("code is required"),
		},
		{
			name: "unsupported comparator",
			modify: func(cn ConditionNode) ConditionNode {
				cn.Comparator = "unsupported"
				return cn
			},
			wantErr: ValidationError("unsupported comparator: unsupported"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cn := tc.modify(baseCond)
			err := cn.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidatePeriod(t *testing.T) {
	tests := []struct {
		name    string
		period  Period
		wantErr error
	}{
		{
			name: "valid year period",
			period: Period{
				Type:  PeriodTypeYear,
				Years: 1,
			},
		},
		{
			name: "year period with zero years",
			period: Period{
				Type:  PeriodTypeYear,
				Years: 0,
			},
			wantErr: ValidationError("year period must have years greater than 0"),
		},
		{
			name: "valid rolling period with days",
			period: Period{
				Type:        PeriodTypeRolling,
				RollingDays: 1,
			},
		},
		{
			name: "valid rolling period with months",
			period: Period{
				Type:          PeriodTypeRolling,
				RollingMonths: 1,
			},
		},
		{
			name: "valid rolling period with years",
			period: Period{
				Type:         PeriodTypeRolling,
				RollingYears: 1,
			},
		},
		{
			name: "rolling period with no positive days/months/years",
			period: Period{
				Type: PeriodTypeRolling,
			},
			wantErr: ValidationError("rolling period must have one days/months/years greater than 0"),
		},
		{
			name: "unknown period type",
			period: Period{
				Type: "unknown",
			},
			wantErr: ValidationError("unknown period type: unknown"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.period.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

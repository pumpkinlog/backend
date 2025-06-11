package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConditionType_Valid(t *testing.T) {
	tests := []struct {
		name string
		ct   ConditionType
		want bool
	}{
		{
			name: "valid string",
			ct:   ConditionTypeString,
			want: true,
		},
		{
			"valid boolean",
			ConditionTypeBoolean,
			true,
		},
		{
			"valid integer",
			ConditionTypeInteger,
			true,
		},
		{
			"valid select",
			ConditionTypeSelect,
			true,
		},
		{
			"valid multi_select",
			ConditionTypeMultiSelect,
			true,
		},
		{
			"invalid type",
			ConditionType("invalid"),
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.ct.Valid()
			require.Equal(t, tc.want, got)
		})
	}
}

func TestValidateCondition(t *testing.T) {
	condition := Condition{
		ID:       Code("VALID_CODE"),
		RegionID: RegionID("GB"),
		Prompt:   "What is your age?",
		Type:     ConditionTypeInteger,
	}

	tests := []struct {
		name    string
		modify  func(c Condition) Condition
		wantErr error
	}{
		{
			name:   "valid condition",
			modify: func(c Condition) Condition { return c },
		},
		{
			name: "invalid ID",
			modify: func(c Condition) Condition {
				c.ID = Code("")
				return c
			},
			wantErr: ValidationError("code is required"),
		},
		{
			name: "invalid RegionID",
			modify: func(c Condition) Condition {
				c.RegionID = RegionID("")
				return c
			},
			wantErr: ValidationError("region ID is required"),
		},
		{
			name: "empty prompt",
			modify: func(c Condition) Condition {
				c.Prompt = ""
				return c
			},
			wantErr: ValidationError("prompt is required"),
		},
		{
			name: "invalid ConditionType",
			modify: func(c Condition) Condition {
				c.Type = ConditionType("invalid")
				return c
			},
			wantErr: ValidationError("invalid condition type: invalid"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.modify(condition)
			err := c.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

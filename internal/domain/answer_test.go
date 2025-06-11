package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateAnswer(t *testing.T) {
	timestamp := time.Now()
	answer := Answer{
		UserID:      123,
		ConditionID: "COND-ID",
		RegionID:    "JE",
		Value:       struct{}{},
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	tests := []struct {
		name    string
		modify  func(a Answer) Answer
		wantErr error
	}{
		{
			name: "valid answer",
			modify: func(a Answer) Answer {
				return a
			},
		},
		{
			name: "missing/invalid userID",
			modify: func(a Answer) Answer {
				a.UserID = -1
				return a
			},
			wantErr: ValidationError("user ID is required"),
		},
		{
			name: "missing conditionID",
			modify: func(a Answer) Answer {
				a.ConditionID = ""
				return a
			},
			wantErr: ValidationError("code is required"),
		},
		{
			name: "long conditionID",
			modify: func(a Answer) Answer {
				a.ConditionID = Code(strings.Repeat("a", 129))
				return a
			},
			wantErr: ValidationError("code must be less than 128 characters"),
		},
		{
			name: "invalid conditionID",
			modify: func(a Answer) Answer {
				a.ConditionID = `'"invalid()code+=`
				return a
			},
			wantErr: ValidationError("code must match regular expression: %s", codeRegex),
		},
		{
			name: "missing regionID",
			modify: func(a Answer) Answer {
				a.RegionID = ""
				return a
			},
			wantErr: ValidationError("region ID is required"),
		},
		{
			name: "invalid regionID",
			modify: func(a Answer) Answer {
				a.RegionID = "zzz-zzz"
				return a
			},
			wantErr: ValidationError("region ID must match regular expression: %s", regionIDRegex),
		},
		{
			name: "missing value",
			modify: func(a Answer) Answer {
				a.Value = nil
				return a
			},
			wantErr: ValidationError("value is required"),
		},
		{
			name: "missing created at",
			modify: func(a Answer) Answer {
				a.CreatedAt = time.Time{}
				return a
			},
			wantErr: ValidationError("created at is required"),
		},
		{
			name: "missing updated at",
			modify: func(a Answer) Answer {
				a.UpdatedAt = time.Time{}
				return a
			},
			wantErr: ValidationError("updated at is required"),
		},
		{
			name: "created at in the future",
			modify: func(a Answer) Answer {
				a.CreatedAt = a.CreatedAt.Add(time.Second)
				return a
			},
			wantErr: ValidationError("created at timestamp cannot be in the future"),
		},
		{
			name: "updated at in the future",
			modify: func(a Answer) Answer {
				a.UpdatedAt = a.UpdatedAt.Add(time.Second)
				return a
			},
			wantErr: ValidationError("updated at timestamp cannot be in the future"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ans := tc.modify(answer)
			err := ans.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

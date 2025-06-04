package comparator

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name       string
		condType   domain.ConditionType
		comparator domain.Comparator
		expected   any
		actual     any
		want       bool
		wantErr    error
	}{
		// Boolean tests
		{
			name:       "boolean equals true",
			condType:   domain.ConditionTypeBoolean,
			comparator: domain.ComparatorEquals,
			expected:   true,
			actual:     true,
			want:       true,
		},
		{
			name:       "boolean equals false",
			condType:   domain.ConditionTypeBoolean,
			comparator: domain.ComparatorEquals,
			expected:   false,
			actual:     false,
			want:       true,
		},
		{
			name:       "boolean not equals",
			condType:   domain.ConditionTypeBoolean,
			comparator: domain.ComparatorNotEquals,
			expected:   true,
			actual:     false,
			want:       true,
		},
		{
			name:       "boolean invalid operator",
			condType:   domain.ConditionTypeBoolean,
			comparator: domain.ComparatorGreaterThan,
			expected:   true,
			actual:     false,
			wantErr:    ErrInvalidComparator,
		},

		// Integer tests
		{
			name:       "integer equals",
			condType:   domain.ConditionTypeInteger,
			comparator: domain.ComparatorEquals,
			expected:   42,
			actual:     42,
			want:       true,
		},
		{
			name:       "integer greater than",
			condType:   domain.ConditionTypeInteger,
			comparator: domain.ComparatorGreaterThan,
			expected:   50,
			actual:     42,
			want:       true,
		},
		{
			name:       "integer less than",
			condType:   domain.ConditionTypeInteger,
			comparator: domain.ComparatorLessThan,
			expected:   30,
			actual:     42,
			want:       true,
		},
		{
			name:       "integer string conversion",
			condType:   domain.ConditionTypeInteger,
			comparator: domain.ComparatorEquals,
			expected:   "42",
			actual:     42,
			want:       true,
		},

		// String tests
		{
			name:       "string equals",
			condType:   domain.ConditionTypeString,
			comparator: domain.ComparatorEquals,
			expected:   "hello",
			actual:     "hello",
			want:       true,
		},
		{
			name:       "string contains",
			condType:   domain.ConditionTypeString,
			comparator: domain.ComparatorContains,
			expected:   "hello world",
			actual:     "world",
			want:       true,
		},
		{
			name:       "string not contains",
			condType:   domain.ConditionTypeString,
			comparator: domain.ComparatorContains,
			expected:   "hello",
			actual:     "world",
			want:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Compare(tc.condType, tc.comparator, tc.expected, tc.actual)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tc.wantErr), "expected error %v, got %v", tc.wantErr, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestValidateOperator(t *testing.T) {
	tests := []struct {
		name     string
		condType domain.ConditionType
		operator domain.Comparator
		wantErr  bool
	}{
		{
			name:     "valid boolean operator",
			condType: domain.ConditionTypeBoolean,
			operator: domain.ComparatorEquals,
			wantErr:  false,
		},
		{
			name:     "invalid boolean operator",
			condType: domain.ConditionTypeBoolean,
			operator: domain.ComparatorGreaterThan,
			wantErr:  true,
		},
		{
			name:     "valid integer operator",
			condType: domain.ConditionTypeInteger,
			operator: domain.ComparatorGreaterThan,
			wantErr:  false,
		},
		{
			name:     "invalid integer operator",
			condType: domain.ConditionTypeInteger,
			operator: domain.ComparatorContains,
			wantErr:  true,
		},
		{
			name:     "valid string operator",
			condType: domain.ConditionTypeString,
			operator: domain.ComparatorContains,
			wantErr:  false,
		},
		{
			name:     "invalid string operator",
			condType: domain.ConditionTypeString,
			operator: domain.ComparatorGreaterThan,
			wantErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateComparator(tc.condType, tc.operator)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTypeConversion(t *testing.T) {
	tests := []struct {
		name   string
		fn     func(any) (any, bool)
		input  any
		want   any
		wantOk bool
	}{
		{
			name:   "bool to bool",
			fn:     func(v any) (any, bool) { b, ok := asBool(v); return b, ok },
			input:  true,
			want:   true,
			wantOk: true,
		},
		{
			name:   "string to bool",
			fn:     func(v any) (any, bool) { b, ok := asBool(v); return b, ok },
			input:  "true",
			want:   true,
			wantOk: true,
		},
		{
			name:   "invalid bool",
			fn:     func(v any) (any, bool) { b, ok := asBool(v); return b, ok },
			input:  "not a bool",
			want:   false,
			wantOk: false,
		},
		{
			name:   "int to int64",
			fn:     func(v any) (any, bool) { i, ok := asInt64(v); return i, ok },
			input:  42,
			want:   int64(42),
			wantOk: true,
		},
		{
			name:   "string to int64",
			fn:     func(v any) (any, bool) { i, ok := asInt64(v); return i, ok },
			input:  "42",
			want:   int64(42),
			wantOk: true,
		},
		{
			name:   "invalid int",
			fn:     func(v any) (any, bool) { i, ok := asInt64(v); return i, ok },
			input:  "not a number",
			want:   int64(0),
			wantOk: false,
		},
		{
			name:   "string to string",
			fn:     func(v any) (any, bool) { s, ok := asString(v); return s, ok },
			input:  "hello",
			want:   "hello",
			wantOk: true,
		},
		{
			name:   "stringer to string",
			fn:     func(v any) (any, bool) { s, ok := asString(v); return s, ok },
			input:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want:   "2024-01-01 00:00:00 +0000 UTC",
			wantOk: true,
		},
		{
			name:   "invalid string",
			fn:     func(v any) (any, bool) { s, ok := asString(v); return s, ok },
			input:  42,
			want:   "",
			wantOk: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.fn(tc.input)
			assert.Equal(t, tc.wantOk, ok)
			if tc.wantOk {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

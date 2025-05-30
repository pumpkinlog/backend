package comparator

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name      string
		valueType Type
		operator  Operator
		expected  any
		actual    any
		want      bool
		wantErr   error
	}{
		// Boolean tests
		{
			name:      "boolean equals true",
			valueType: TypeBoolean,
			operator:  OperatorEquals,
			expected:  true,
			actual:    true,
			want:      true,
		},
		{
			name:      "boolean equals false",
			valueType: TypeBoolean,
			operator:  OperatorEquals,
			expected:  false,
			actual:    false,
			want:      true,
		},
		{
			name:      "boolean not equals",
			valueType: TypeBoolean,
			operator:  OperatorNotEquals,
			expected:  true,
			actual:    false,
			want:      true,
		},
		{
			name:      "boolean invalid operator",
			valueType: TypeBoolean,
			operator:  OperatorGreaterThan,
			expected:  true,
			actual:    false,
			wantErr:   ErrInvalidOperator,
		},

		// Integer tests
		{
			name:      "integer equals",
			valueType: TypeInteger,
			operator:  OperatorEquals,
			expected:  42,
			actual:    42,
			want:      true,
		},
		{
			name:      "integer greater than",
			valueType: TypeInteger,
			operator:  OperatorGreaterThan,
			expected:  50,
			actual:    42,
			want:      true,
		},
		{
			name:      "integer less than",
			valueType: TypeInteger,
			operator:  OperatorLessThan,
			expected:  30,
			actual:    42,
			want:      true,
		},
		{
			name:      "integer string conversion",
			valueType: TypeInteger,
			operator:  OperatorEquals,
			expected:  "42",
			actual:    42,
			want:      true,
		},

		// String tests
		{
			name:      "string equals",
			valueType: TypeString,
			operator:  OperatorEquals,
			expected:  "hello",
			actual:    "hello",
			want:      true,
		},
		{
			name:      "string contains",
			valueType: TypeString,
			operator:  OperatorContains,
			expected:  "hello world",
			actual:    "world",
			want:      true,
		},
		{
			name:      "string not contains",
			valueType: TypeString,
			operator:  OperatorContains,
			expected:  "hello",
			actual:    "world",
			want:      false,
		},

		// Date tests
		{
			name:      "date equals",
			valueType: TypeDate,
			operator:  OperatorEquals,
			expected:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			actual:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      true,
		},
		{
			name:      "date greater than",
			valueType: TypeDate,
			operator:  OperatorGreaterThan,
			expected:  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			actual:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      true,
		},
		{
			name:      "date string conversion",
			valueType: TypeDate,
			operator:  OperatorEquals,
			expected:  "2024-01-01T00:00:00Z",
			actual:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      true,
		},

		// Float tests
		{
			name:      "float equals",
			valueType: TypeFloat,
			operator:  OperatorEquals,
			expected:  3.14,
			actual:    3.14,
			want:      true,
		},
		{
			name:      "float greater than",
			valueType: TypeFloat,
			operator:  OperatorGreaterThan,
			expected:  4.0,
			actual:    3.14,
			want:      true,
		},
		{
			name:      "float string conversion",
			valueType: TypeFloat,
			operator:  OperatorEquals,
			expected:  "3.14",
			actual:    3.14,
			want:      true,
		},

		// Error cases
		{
			name:      "unsupported type",
			valueType: "unsupported",
			operator:  OperatorEquals,
			expected:  "value",
			actual:    "value",
			wantErr:   ErrUnsupportedType,
		},
		{
			name:      "type conversion failed",
			valueType: TypeInteger,
			operator:  OperatorEquals,
			expected:  "not a number",
			actual:    42,
			wantErr:   ErrConversionFailed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Compare(tc.valueType, tc.operator, tc.expected, tc.actual)

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
		name      string
		valueType Type
		operator  Operator
		wantErr   bool
	}{
		{
			name:      "valid boolean operator",
			valueType: TypeBoolean,
			operator:  OperatorEquals,
			wantErr:   false,
		},
		{
			name:      "invalid boolean operator",
			valueType: TypeBoolean,
			operator:  OperatorGreaterThan,
			wantErr:   true,
		},
		{
			name:      "valid integer operator",
			valueType: TypeInteger,
			operator:  OperatorGreaterThan,
			wantErr:   false,
		},
		{
			name:      "invalid integer operator",
			valueType: TypeInteger,
			operator:  OperatorContains,
			wantErr:   true,
		},
		{
			name:      "valid string operator",
			valueType: TypeString,
			operator:  OperatorContains,
			wantErr:   false,
		},
		{
			name:      "invalid string operator",
			valueType: TypeString,
			operator:  OperatorGreaterThan,
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateOperator(tc.valueType, tc.operator)
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
		// Boolean conversion
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

		// Integer conversion
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

		// String conversion
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

		// Date conversion
		{
			name:   "time to time",
			fn:     func(v any) (any, bool) { t, ok := asTime(v); return t, ok },
			input:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantOk: true,
		},
		{
			name:   "string to time",
			fn:     func(v any) (any, bool) { t, ok := asTime(v); return t, ok },
			input:  "2024-01-01T00:00:00Z",
			want:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantOk: true,
		},
		{
			name:   "invalid time",
			fn:     func(v any) (any, bool) { t, ok := asTime(v); return t, ok },
			input:  "not a date",
			want:   time.Time{},
			wantOk: false,
		},

		// Float conversion
		{
			name:   "float to float64",
			fn:     func(v any) (any, bool) { f, ok := asFloat64(v); return f, ok },
			input:  3.14,
			want:   float64(3.14),
			wantOk: true,
		},
		{
			name:   "string to float64",
			fn:     func(v any) (any, bool) { f, ok := asFloat64(v); return f, ok },
			input:  "3.14",
			want:   float64(3.14),
			wantOk: true,
		},
		{
			name:   "invalid float",
			fn:     func(v any) (any, bool) { f, ok := asFloat64(v); return f, ok },
			input:  "not a number",
			want:   float64(0),
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

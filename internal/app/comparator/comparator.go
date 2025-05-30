package comparator

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Type represents the type of value being compared
type Type string

const (
	TypeBoolean Type = "boolean"
	TypeInteger Type = "integer"
	TypeString  Type = "string"
	TypeDate    Type = "date"
	TypeFloat   Type = "float"
)

// Operator represents the comparison operator
type Operator string

const (
	OperatorEquals            Operator = "eq"
	OperatorNotEquals         Operator = "neq"
	OperatorGreaterThan       Operator = "gt"
	OperatorGreaterThanEquals Operator = "gte"
	OperatorLessThan          Operator = "lt"
	OperatorLessThanEquals    Operator = "lte"
	OperatorContains          Operator = "contains"
	OperatorIn                Operator = "in"
	OperatorNotIn             Operator = "not_in"
)

// Error types for comparison operations
var (
	ErrTypeMismatch     = fmt.Errorf("type mismatch")
	ErrInvalidOperator  = fmt.Errorf("invalid operator")
	ErrInvalidValue     = fmt.Errorf("invalid value")
	ErrUnsupportedType  = fmt.Errorf("unsupported type")
	ErrConversionFailed = fmt.Errorf("conversion failed")
)

// Compare compares two values using the specified operator
func Compare(valueType Type, operator Operator, expected, actual any) (bool, error) {

	if err := validateOperator(valueType, operator); err != nil {
		return false, fmt.Errorf("validate operator: %w", err)
	}

	exp, act, err := convertValues(valueType, expected, actual)
	if err != nil {
		return false, fmt.Errorf("convert values: %w", err)
	}

	switch valueType {
	case TypeBoolean:
		return compareBools(exp.(bool), act.(bool), operator)
	case TypeInteger:
		return compareInts(exp.(int64), act.(int64), operator)
	case TypeString:
		return compareStrings(exp.(string), act.(string), operator)
	case TypeDate:
		return compareDates(exp.(time.Time), act.(time.Time), operator)
	case TypeFloat:
		return compareFloats(exp.(float64), act.(float64), operator)
	default:
		return false, fmt.Errorf("%w: %s", ErrUnsupportedType, valueType)
	}
}

var validOperators = map[Type][]Operator{
	TypeBoolean: {
		OperatorEquals,
		OperatorNotEquals,
	},
	TypeInteger: {
		OperatorEquals,
		OperatorNotEquals,
		OperatorGreaterThan,
		OperatorGreaterThanEquals,
		OperatorLessThan,
		OperatorLessThanEquals,
	},
	TypeString: {
		OperatorEquals,
		OperatorNotEquals,
		OperatorContains,
		OperatorIn,
		OperatorNotIn,
	},
	TypeDate: {
		OperatorEquals,
		OperatorNotEquals,
		OperatorGreaterThan,
		OperatorGreaterThanEquals,
		OperatorLessThan,
		OperatorLessThanEquals,
	},
	TypeFloat: {
		OperatorEquals,
		OperatorNotEquals,
		OperatorGreaterThan,
		OperatorGreaterThanEquals,
		OperatorLessThan,
		OperatorLessThanEquals,
	},
}

// validateOperator checks if the operator is valid for the given type
func validateOperator(valueType Type, operator Operator) error {

	operators, ok := validOperators[valueType]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedType, valueType)
	}

	if !slices.Contains(operators, operator) {
		return fmt.Errorf("%w: %s for type %s", ErrInvalidOperator, operator, valueType)
	}

	return nil
}

// convertValues converts the input values to their proper types
func convertValues(valueType Type, expected, actual any) (any, any, error) {
	switch valueType {
	case TypeBoolean:
		exp, ok1 := asBool(expected)
		act, ok2 := asBool(actual)
		if !ok1 || !ok2 {
			return nil, nil, fmt.Errorf("%w: boolean conversion failed", ErrConversionFailed)
		}
		return exp, act, nil

	case TypeInteger:
		exp, ok1 := asInt64(expected)
		act, ok2 := asInt64(actual)
		if !ok1 || !ok2 {
			return nil, nil, fmt.Errorf("%w: integer conversion failed", ErrConversionFailed)
		}
		return exp, act, nil

	case TypeString:
		exp, ok1 := asString(expected)
		act, ok2 := asString(actual)
		if !ok1 || !ok2 {
			return nil, nil, fmt.Errorf("%w: string conversion failed", ErrConversionFailed)
		}
		return exp, act, nil

	case TypeDate:
		exp, ok1 := asTime(expected)
		act, ok2 := asTime(actual)
		if !ok1 || !ok2 {
			return nil, nil, fmt.Errorf("%w: date conversion failed", ErrConversionFailed)
		}
		return exp, act, nil

	case TypeFloat:
		exp, ok1 := asFloat64(expected)
		act, ok2 := asFloat64(actual)
		if !ok1 || !ok2 {
			return nil, nil, fmt.Errorf("%w: float conversion failed", ErrConversionFailed)
		}
		return exp, act, nil

	default:
		return nil, nil, fmt.Errorf("%w: %s", ErrUnsupportedType, valueType)
	}
}

func asBool(v any) (bool, bool) {
	switch val := v.(type) {
	case bool:
		return val, true
	case string:
		b, err := strconv.ParseBool(val)
		return b, err == nil
	default:
		return false, false
	}
}

func asInt64(v any) (int64, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int8:
		return int64(val), true
	case int16:
		return int64(val), true
	case int32:
		return int64(val), true
	case int64:
		return val, true
	case string:
		i, err := strconv.ParseInt(val, 10, 64)
		return i, err == nil
	default:
		return 0, false
	}
}

func asString(v any) (string, bool) {
	switch val := v.(type) {
	case string:
		return val, true
	case fmt.Stringer:
		return val.String(), true
	default:
		return "", false
	}
}

func asTime(v any) (time.Time, bool) {
	switch val := v.(type) {
	case time.Time:
		return val, true
	case string:
		t, err := time.Parse(time.RFC3339, val)
		return t, err == nil
	default:
		return time.Time{}, false
	}
}

func asFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case string:
		f, err := strconv.ParseFloat(val, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func compareBools(a, b bool, op Operator) (bool, error) {
	switch op {
	case OperatorEquals:
		return a == b, nil
	case OperatorNotEquals:
		return a != b, nil
	default:
		return false, fmt.Errorf("%w: %s for boolean", ErrInvalidOperator, op)
	}
}

func compareInts(a, b int64, op Operator) (bool, error) {
	switch op {
	case OperatorEquals:
		return a == b, nil
	case OperatorNotEquals:
		return a != b, nil
	case OperatorGreaterThan:
		return a > b, nil
	case OperatorGreaterThanEquals:
		return a >= b, nil
	case OperatorLessThan:
		return a < b, nil
	case OperatorLessThanEquals:
		return a <= b, nil
	default:
		return false, fmt.Errorf("%w: %s for integer", ErrInvalidOperator, op)
	}
}

func compareStrings(a, b string, op Operator) (bool, error) {
	switch op {
	case OperatorEquals:
		return a == b, nil
	case OperatorNotEquals:
		return a != b, nil
	case OperatorContains:
		return strings.Contains(a, b), nil
	case OperatorIn:
		return strings.Contains(a, b), nil
	case OperatorNotIn:
		return !strings.Contains(a, b), nil
	default:
		return false, fmt.Errorf("%w: %s for string", ErrInvalidOperator, op)
	}
}

func compareDates(a, b time.Time, op Operator) (bool, error) {
	switch op {
	case OperatorEquals:
		return a.Equal(b), nil
	case OperatorNotEquals:
		return !a.Equal(b), nil
	case OperatorGreaterThan:
		return a.After(b), nil
	case OperatorGreaterThanEquals:
		return a.After(b) || a.Equal(b), nil
	case OperatorLessThan:
		return a.Before(b), nil
	case OperatorLessThanEquals:
		return a.Before(b) || a.Equal(b), nil
	default:
		return false, fmt.Errorf("%w: %s for date", ErrInvalidOperator, op)
	}
}

func compareFloats(a, b float64, op Operator) (bool, error) {
	switch op {
	case OperatorEquals:
		return a == b, nil
	case OperatorNotEquals:
		return a != b, nil
	case OperatorGreaterThan:
		return a > b, nil
	case OperatorGreaterThanEquals:
		return a >= b, nil
	case OperatorLessThan:
		return a < b, nil
	case OperatorLessThanEquals:
		return a <= b, nil
	default:
		return false, fmt.Errorf("%w: %s for float", ErrInvalidOperator, op)
	}
}

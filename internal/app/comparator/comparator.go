package comparator

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pumpkinlog/backend/internal/domain"
)

// Error types for comparison operations
var (
	ErrTypeMismatch      = fmt.Errorf("type mismatch")
	ErrInvalidComparator = fmt.Errorf("invalid comparator")
	ErrInvalidValue      = fmt.Errorf("invalid value")
	ErrUnsupportedType   = fmt.Errorf("unsupported type")
	ErrConversionFailed  = fmt.Errorf("conversion failed")
)

// Compare compares two values using the specified domain.Comparator
func Compare(condType domain.ConditionType, comparator domain.Comparator, expected, actual any) (bool, error) {

	if err := validateComparator(condType, comparator); err != nil {
		return false, fmt.Errorf("validate domain.Comparator: %w", err)
	}

	exp, act, err := convertValues(condType, expected, actual)
	if err != nil {
		return false, fmt.Errorf("convert values: %w", err)
	}

	switch condType {
	case domain.ConditionTypeBoolean:
		return compareBools(exp.(bool), act.(bool), comparator)
	case domain.ConditionTypeInteger:
		return compareInts(exp.(int64), act.(int64), comparator)
	case domain.ConditionTypeString:
		return compareStrings(exp.(string), act.(string), comparator)
	default:
		return false, fmt.Errorf("%w: %s", ErrUnsupportedType, condType)
	}
}

var validComparators = map[domain.ConditionType][]domain.Comparator{
	domain.ConditionTypeBoolean: {
		domain.ComparatorEquals,
		domain.ComparatorNotEquals,
	},
	domain.ConditionTypeInteger: {
		domain.ComparatorEquals,
		domain.ComparatorNotEquals,
		domain.ComparatorGreaterThan,
		domain.ComparatorGreaterThanOrEquals,
		domain.ComparatorLessThan,
		domain.ComparatorLessThanOrEquals,
	},
	domain.ConditionTypeString: {
		domain.ComparatorEquals,
		domain.ComparatorNotEquals,
		domain.ComparatorContains,
		domain.ComparatorIn,
		domain.ComparatorNotIn,
	},
}

// validateComparator checks if the comparator is valid for the given type
func validateComparator(condType domain.ConditionType, comparator domain.Comparator) error {

	comparators, ok := validComparators[condType]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedType, condType)
	}

	if !slices.Contains(comparators, comparator) {
		return fmt.Errorf("%w: %s for type %s", ErrInvalidComparator, comparator, condType)
	}

	return nil
}

// convertValues converts the input values to their proper types
func convertValues(condType domain.ConditionType, expected, actual any) (any, any, error) {
	switch condType {
	case domain.ConditionTypeBoolean:
		exp, ok1 := asBool(expected)
		act, ok2 := asBool(actual)
		if !ok1 || !ok2 {
			return nil, nil, fmt.Errorf("%w: boolean conversion failed", ErrConversionFailed)
		}
		return exp, act, nil

	case domain.ConditionTypeInteger:
		exp, ok1 := asInt64(expected)
		act, ok2 := asInt64(actual)
		if !ok1 || !ok2 {
			return nil, nil, fmt.Errorf("%w: integer conversion failed", ErrConversionFailed)
		}
		return exp, act, nil

	case domain.ConditionTypeString:
		exp, ok1 := asString(expected)
		act, ok2 := asString(actual)
		if !ok1 || !ok2 {
			return nil, nil, fmt.Errorf("%w: string conversion failed", ErrConversionFailed)
		}
		return exp, act, nil

	default:
		return nil, nil, fmt.Errorf("%w: %s", ErrUnsupportedType, condType)
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

func compareBools(a, b bool, cmp domain.Comparator) (bool, error) {
	switch cmp {
	case domain.ComparatorEquals:
		return a == b, nil
	case domain.ComparatorNotEquals:
		return a != b, nil
	default:
		return false, fmt.Errorf("%w: %s for boolean", ErrInvalidComparator, cmp)
	}
}

func compareInts(a, b int64, cmp domain.Comparator) (bool, error) {
	switch cmp {
	case domain.ComparatorEquals:
		return a == b, nil
	case domain.ComparatorNotEquals:
		return a != b, nil
	case domain.ComparatorGreaterThan:
		return a > b, nil
	case domain.ComparatorGreaterThanOrEquals:
		return a >= b, nil
	case domain.ComparatorLessThan:
		return a < b, nil
	case domain.ComparatorLessThanOrEquals:
		return a <= b, nil
	default:
		return false, fmt.Errorf("%w: %s for integer", ErrInvalidComparator, cmp)
	}
}

func compareStrings(a, b string, op domain.Comparator) (bool, error) {
	switch op {
	case domain.ComparatorEquals:
		return a == b, nil
	case domain.ComparatorNotEquals:
		return a != b, nil
	case domain.ComparatorContains:
		return strings.Contains(a, b), nil
	case domain.ComparatorIn:
		return strings.Contains(a, b), nil
	case domain.ComparatorNotIn:
		return !strings.Contains(a, b), nil
	default:
		return false, fmt.Errorf("%w: %s for string", ErrInvalidComparator, op)
	}
}

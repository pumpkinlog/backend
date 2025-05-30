package period

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
)

func TestComputeRulePeriod(t *testing.T) {

	tests := []struct {
		name     string
		asOf     time.Time
		rule     *domain.Rule
		expected domain.TimeWindow
	}{
		{
			name: "calendar year period",
			asOf: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			rule: &domain.Rule{
				PeriodType:     domain.PeriodTypeYear,
				YearStartMonth: time.January,
				YearStartDay:   1,
				Years:          1,
			},
			expected: domain.TimeWindow{
				Start: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name: "four calendar year period",
			asOf: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			rule: &domain.Rule{
				PeriodType:     domain.PeriodTypeYear,
				YearStartMonth: time.January,
				YearStartDay:   1,
				Years:          4,
			},
			expected: domain.TimeWindow{
				Start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name: "fiscal year period",
			asOf: time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
			rule: &domain.Rule{
				PeriodType:     domain.PeriodTypeYear,
				YearStartMonth: time.July,
				YearStartDay:   1,
				Years:          1,
			},
			expected: domain.TimeWindow{
				Start: time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2024, time.June, 30, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name: "four fiscal year period",
			asOf: time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
			rule: &domain.Rule{
				PeriodType:     domain.PeriodTypeYear,
				YearStartMonth: time.July,
				YearStartDay:   1,
				Years:          4,
			},
			expected: domain.TimeWindow{
				Start: time.Date(2022, time.July, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.June, 30, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name: "rolling period",
			asOf: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			rule: &domain.Rule{
				PeriodType:     domain.PeriodTypeRolling,
				YearStartMonth: time.January,
				YearStartDay:   1,
				RollingDays:    90,
			},
			expected: domain.TimeWindow{
				Start: time.Date(2024, time.October, 3, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ComputeRulePeriod(tc.asOf, tc.rule)
			if result.Start != tc.expected.Start || result.End != tc.expected.End {
				require.Equal(t, tc.expected.Start, result.Start)
			}
		})
	}
}

func TestComputeRulesPeriod(t *testing.T) {

	tests := []struct {
		name     string
		asOf     time.Time
		rules    []*domain.Rule
		expected domain.TimeWindow
	}{
		{
			name: "multiple rules with calendar year period",
			asOf: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			rules: []*domain.Rule{
				{
					RegionID:       "JE",
					PeriodType:     domain.PeriodTypeYear,
					YearStartMonth: time.January,
					YearStartDay:   1,
					Years:          1,
				},
				{
					RegionID:       "JE",
					PeriodType:     domain.PeriodTypeYear,
					YearStartMonth: time.January,
					YearStartDay:   1,
					Years:          4,
				},
			},
			expected: domain.TimeWindow{
				Start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ComputeRulesPeriod(tc.asOf, tc.rules)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestComputePeriodByRegion(t *testing.T) {

	tests := []struct {
		name     string
		asOf     time.Time
		rules    []*domain.Rule
		expected map[string]domain.TimeWindow
	}{
		{
			name: "single rules with calendar year period",
			asOf: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			rules: []*domain.Rule{
				{
					RegionID:       "JE",
					PeriodType:     domain.PeriodTypeYear,
					YearStartMonth: time.January,
					YearStartDay:   1,
					Years:          1,
				},
			},
			expected: map[string]domain.TimeWindow{
				"JE": {
					Start: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
				},
			},
		},
		{
			name: "multiple rules with different periods",
			asOf: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			rules: []*domain.Rule{
				{
					RegionID:       "JE",
					PeriodType:     domain.PeriodTypeYear,
					YearStartMonth: time.January,
					YearStartDay:   1,
					Years:          1,
				},
				{
					RegionID:       "JE",
					PeriodType:     domain.PeriodTypeYear,
					YearStartMonth: time.January,
					YearStartDay:   1,
					Years:          4,
				},
			},
			expected: map[string]domain.TimeWindow{
				"JE": {
					Start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ComputePeriodByRegion(tc.asOf, tc.rules)
			require.Equal(t, tc.expected, result)
		})
	}
}

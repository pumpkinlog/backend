package period

import (
	"testing"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestComputeYearPeriod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		asOf        time.Time
		startMonth  time.Month
		startDay    int
		offsetYears int
		years       int
		expected    domain.TimeWindow
	}{
		{
			name:       "Current calendar year",
			asOf:       time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			startMonth: time.January,
			startDay:   1,
			expected: domain.TimeWindow{
				Start: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:        "Previous calendar year",
			asOf:        time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			startMonth:  time.January,
			startDay:    1,
			offsetYears: 1,
			expected: domain.TimeWindow{
				Start: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2024, time.December, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:        "Next calendar year",
			asOf:        time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			startMonth:  time.January,
			startDay:    1,
			offsetYears: -1,
			expected: domain.TimeWindow{
				Start: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2026, time.December, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:       "Current fiscal year",
			asOf:       time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			startMonth: time.July,
			startDay:   1,
			expected: domain.TimeWindow{
				Start: time.Date(2024, time.July, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.June, 30, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:        "Previous fiscal year",
			asOf:        time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			startMonth:  time.July,
			startDay:    1,
			offsetYears: 1,
			expected: domain.TimeWindow{
				Start: time.Date(2023, time.July, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2024, time.June, 30, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:       "Previous four calendar years",
			asOf:       time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			startMonth: time.January,
			startDay:   1,
			years:      4,
			expected: domain.TimeWindow{
				Start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:       "Previous four fiscal years",
			asOf:       time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			startMonth: time.July,
			startDay:   1,
			years:      4,
			expected: domain.TimeWindow{
				Start: time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.June, 30, 23, 59, 59, 0, time.UTC),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			period := computeYearPeriod(tc.asOf, tc.startMonth, tc.startDay, tc.years, tc.offsetYears)

			require.Equal(t, tc.expected, period)
		})
	}
}

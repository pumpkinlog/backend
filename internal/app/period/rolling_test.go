package period

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
)

func TestComputeRollingPeriod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		asOf     time.Time
		years    int
		months   int
		days     int
		expected domain.TimeWindow
	}{
		{
			name:   "Rolling 90 days",
			asOf:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			years:  0,
			months: 0,
			days:   90,
			expected: domain.TimeWindow{
				Start: time.Date(2024, time.October, 3, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.January, 1, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:   "Rolling 1 year",
			asOf:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			years:  1,
			months: 0,
			days:   0,
			expected: domain.TimeWindow{
				Start: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.January, 1, 23, 59, 59, 0, time.UTC),
			},
		},
		{
			name:   "Rolling 365 days",
			asOf:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			years:  0,
			months: 0,
			days:   365,
			expected: domain.TimeWindow{
				Start: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2025, time.January, 1, 23, 59, 59, 0, time.UTC),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			period := computeRollingPeriod(tc.asOf, tc.years, tc.months, tc.days)

			require.Equal(t, period, tc.expected)
		})
	}
}

package period

import (
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

func computeRollingPeriod(asOf time.Time, years, months, days int) domain.TimeWindow {
	start := asOf.AddDate(-years, -months, -days)
	end := asOf.Truncate(24*time.Hour).AddDate(0, 0, 1).Add(-time.Second)

	return domain.TimeWindow{
		Start: start,
		End:   end,
	}
}

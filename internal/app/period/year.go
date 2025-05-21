package period

import (
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

func computeYearPeriod(asOf time.Time, startMonth time.Month, startDay, years, offsetYears int) domain.TimeWindow {

	year := asOf.Year()
	boundary := time.Date(year, startMonth, startDay, 0, 0, 0, 0, asOf.Location())

	if asOf.Before(boundary) {
		year--
	}

	if years <= 0 {
		years = 1
	}

	finalYear := year - offsetYears
	startYear := finalYear - (years - 1)

	start := time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, asOf.Location())
	end := time.Date(finalYear, startMonth, startDay, 0, 0, 0, 0, asOf.Location()).AddDate(1, 0, 0).Add(-time.Second)

	return domain.TimeWindow{Start: start, End: end}
}

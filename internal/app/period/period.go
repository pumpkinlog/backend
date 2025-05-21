package period

import (
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

func ComputePeriod(asOf time.Time, rule *domain.Rule) domain.TimeWindow {
	switch rule.PeriodType {
	case domain.PeriodTypeYear:
		return computeYearPeriod(asOf, rule.YearStartMonth, rule.YearStartDay, rule.Years, rule.OffsetYears)
	case domain.PeriodTypeRolling:
		return computeRollingPeriod(asOf, rule.RollingYears, rule.RollingMonths, rule.RollingDays)
	default:
		return domain.TimeWindow{
			Start: time.Now().UTC(),
			End:   time.Now().UTC(),
		}
	}
}

func ComputePeriodByRegion(asOf time.Time, rules []*domain.Rule) map[string]domain.TimeWindow {
	bounds := make(map[string]domain.TimeWindow, len(rules))

	for _, rule := range rules {
		bounds[rule.RegionID] = ComputePeriod(asOf, rule)
	}

	return bounds
}

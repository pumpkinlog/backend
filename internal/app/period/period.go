package period

import (
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

func ComputeRulePeriod(asOf time.Time, rule *domain.Rule) domain.TimeWindow {
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

func ComputeRulesPeriod(asOf time.Time, rules []*domain.Rule) domain.TimeWindow {
	var bounds domain.TimeWindow

	for i, rule := range rules {
		b := ComputeRulePeriod(asOf, rule)

		if i == 0 {
			bounds = b
			continue
		}

		if b.Start.Before(bounds.Start) {
			bounds.Start = b.Start
		}

		if b.End.After(bounds.End) {
			bounds.End = b.End
		}
	}

	return bounds
}

func ComputePeriodByRegion(asOf time.Time, rules []*domain.Rule) map[string]domain.TimeWindow {
	bounds := make(map[string]domain.TimeWindow, len(rules))

	for _, rule := range rules {
		bounds[rule.RegionID] = ComputeRulePeriod(asOf, rule)
	}

	return bounds
}

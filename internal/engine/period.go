package engine

import (
	"fmt"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

func ComputePeriod(at time.Time, region *domain.Region, period domain.Period) (time.Time, time.Time, error) {
	switch period.Type {
	case "year":
		year := at.Year()
		boundary := time.Date(year, region.YearStartMonth, region.YearStartDay, 0, 0, 0, 0, at.Location())

		if at.Before(boundary) {
			year--
		}

		years := period.Years
		if years <= 0 {
			years = 1
		}

		finalYear := year - period.OffsetYears
		startYear := finalYear - (years - 1)

		start := time.Date(startYear, region.YearStartMonth, region.YearStartDay, 0, 0, 0, 0, at.Location())
		end := time.Date(finalYear, region.YearStartMonth, region.YearStartDay, 0, 0, 0, 0, at.Location()).AddDate(1, 0, 0).Add(-time.Second)

		return start, end, nil

	case "rolling":
		start := at.AddDate(-period.RollingYears, -period.RollingMonths, -period.RollingDays)
		end := at.Truncate(24*time.Hour).AddDate(0, 0, 1).Add(-time.Second)
		return start, end, nil

	default:
		return time.Time{}, time.Time{}, fmt.Errorf("unsupported period type: %s", period.Type)
	}
}

func ComputeMaxPeriod(at time.Time, region *domain.Region, rules []*domain.Rule) (time.Time, time.Time, error) {
	var minStart, maxEnd time.Time

	for _, rule := range rules {
		if rule.Node.Type != domain.NodeTypeStrategy {
			continue
		}

		var sn domain.EvaluatorNode
		if err := rule.Node.Unmarshal(&sn); err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("cannot unmarshal strategy node for rule %s: %w", rule.ID, err)
		}

		start, end, err := ComputePeriod(at, region, sn.Period)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("compute period for rule %s: %w", rule.ID, err)
		}

		if minStart.IsZero() || start.Before(minStart) {
			minStart = start
		}

		if end.After(maxEnd) {
			maxEnd = end
		}
	}

	return minStart, maxEnd, nil
}

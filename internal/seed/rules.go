package seed

import (
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (s *Seeder) rules() []*domain.Rule {
	return []*domain.Rule{
		{
			ID:             RuleJE183Day,
			RegionID:       "JE",
			Name:           "183 Day Rule",
			Description:    "You are present for 183 days or more in a calendar year.",
			RuleType:       domain.RuleTypeAggregate,
			PeriodType:     domain.PeriodTypeYear,
			Threshold:      183,
			YearStartMonth: time.January,
			YearStartDay:   1,
		},
		{
			ID:             RuleJEAbodeOneNight,
			RegionID:       "JE",
			Name:           "Place of abode with one night",
			Description:    "Maintains a place of abode in Jersey and stays one night in a calendar year.",
			RuleType:       domain.RuleTypeAggregate,
			PeriodType:     domain.PeriodTypeYear,
			Threshold:      1,
			YearStartMonth: time.January,
			YearStartDay:   1,
		},
		{
			ID:             RuleJEAvgPresence4Yr,
			RegionID:       "JE",
			Name:           "Average presence over 4 years",
			Description:    "Stays for an average of 3 months per year over 4 years without a place of abode.",
			RuleType:       domain.RuleTypeAverage,
			PeriodType:     domain.PeriodTypeYear,
			Threshold:      1,
			Years:          4,
			YearStartMonth: time.January,
			YearStartDay:   1,
		},
		{
			ID:             RuleGG182Day,
			RegionID:       "GG",
			Name:           "182 Day Rule",
			Description:    "You are present for 182 days or more in a calendar year.",
			RuleType:       domain.RuleTypeAggregate,
			PeriodType:     domain.PeriodTypeYear,
			Threshold:      182,
			YearStartMonth: time.January,
			YearStartDay:   1,
		},
	}
}

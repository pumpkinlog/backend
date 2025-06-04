package seed

import (
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

const (
	Rule183Day = iota
	Rule182Day
	RuleAbodeOneNight
	RuleAvgPresence4Yr
)

var rules = []domain.Rule{
	{
		ID:             Rule182Day,
		RegionID:       "GG",
		Name:           "182 Day Rule",
		Description:    "You are present for 182 days or more in a tax year.",
		RuleType:       domain.RuleTypeAggregate,
		PeriodType:     domain.PeriodTypeYear,
		Threshold:      182,
		YearStartMonth: time.January,
		YearStartDay:   1,
	},
	{
		ID:             Rule183Day,
		RegionID:       "JE",
		Name:           "183 Day Rule",
		Description:    "You are present for 183 days or more in a tax year.",
		RuleType:       domain.RuleTypeAggregate,
		PeriodType:     domain.PeriodTypeYear,
		Threshold:      183,
		YearStartMonth: time.January,
		YearStartDay:   1,
	},
	{
		ID:             RuleAbodeOneNight,
		RegionID:       "JE",
		Name:           "Place of abode with one night",
		Description:    "You maintain a place of abode and stay one night in a tax year.",
		RuleType:       domain.RuleTypeAggregate,
		PeriodType:     domain.PeriodTypeYear,
		Threshold:      1,
		YearStartMonth: time.January,
		YearStartDay:   1,
	},
	{
		ID:             RuleAvgPresence4Yr,
		RegionID:       "JE",
		Name:           "Average presence over 4 years",
		Description:    "You stay for an average of 3 months per year over 4 years without a place of abode.",
		RuleType:       domain.RuleTypeAverage,
		PeriodType:     domain.PeriodTypeYear,
		Threshold:      1,
		Years:          4,
		YearStartMonth: time.January,
		YearStartDay:   1,
	},
}

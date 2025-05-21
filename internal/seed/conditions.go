package seed

import (
	"github.com/pumpkinlog/backend/internal/domain"
)

func (s *Seeder) conditions() []*domain.Condition {
	return []*domain.Condition{
		{
			ID:         RuleJEAbodeOneNightConditionMaintainAbode,
			RuleID:     RuleJEAbodeOneNight,
			Prompt:     "Do you maintain a place of abode in Jersey?",
			Type:       domain.ConditionTypeBoolean,
			Comparator: domain.ComparatorEquals,
			Expected:   true,
		},
		{
			ID:         RuleJEAvgPresence4YrConditionMaintainAbode,
			RuleID:     RuleJEAvgPresence4Yr,
			Prompt:     "Do you maintain a place of abode in Jersey?",
			Type:       domain.ConditionTypeBoolean,
			Comparator: domain.ComparatorEquals,
			Expected:   false,
		},
	}
}

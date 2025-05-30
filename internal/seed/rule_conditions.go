package seed

import "github.com/pumpkinlog/backend/internal/domain"

func (s *Seeder) ruleConditions() []*domain.RuleCondition {
	return []*domain.RuleCondition{
		{
			RuleID:      RuleJEAbodeOneNight,
			ConditionID: ConditionMaintainAbode,
			Comparator:  domain.ComparatorEquals,
			Expected:    true,
		},
		{
			RuleID:      RuleJEAvgPresence4Yr,
			ConditionID: ConditionMaintainAbode,
			Comparator:  domain.ComparatorEquals,
			Expected:    false,
		},
	}
}

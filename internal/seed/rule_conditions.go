package seed

import "github.com/pumpkinlog/backend/internal/domain"

var ruleConditions = []domain.RuleCondition{
	{
		RuleID:      RuleAbodeOneNight,
		ConditionID: ConditionMaintainAbode,
		RegionID:    "JE",
		Comparator:  domain.ComparatorEquals,
		Expected:    true,
	},
	{
		RuleID:      RuleAvgPresence4Yr,
		ConditionID: ConditionMaintainAbode,
		RegionID:    "JE",
		Comparator:  domain.ComparatorEquals,
		Expected:    false,
	},
}

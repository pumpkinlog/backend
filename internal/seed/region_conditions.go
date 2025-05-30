package seed

import "github.com/pumpkinlog/backend/internal/domain"

func (s *Seeder) regionConditions() []*domain.RegionCondition {
	return []*domain.RegionCondition{
		{
			RegionID:    "JE",
			ConditionID: ConditionMaintainAbode,
		},
	}
}

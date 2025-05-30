package seed

import (
	"github.com/pumpkinlog/backend/internal/domain"
)

func (s *Seeder) conditions() []*domain.Condition {
	return []*domain.Condition{
		{
			ID:     ConditionMaintainAbode,
			Prompt: "Do you maintain a place of abode in this region?",
			Type:   domain.ConditionTypeBoolean,
		},
	}
}

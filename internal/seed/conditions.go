package seed

import (
	"github.com/pumpkinlog/backend/internal/domain"
)

const (
	ConditionMaintainAbode = iota
)

var conditions = []domain.Condition{
	{
		ID:       ConditionMaintainAbode,
		RegionID: "JE",
		Prompt:   "Do you maintain a place of abode in Jersey?",
		Type:     domain.ConditionTypeBoolean,
	},
}

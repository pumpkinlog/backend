package domain

type RegionAggregate struct {
	Region         *Region
	Presences      []*Presence
	Rules          []*Rule
	RuleConditions []*RuleCondition
	Conditions     map[int64]*Condition
	Answers        map[int64]*Answer
}

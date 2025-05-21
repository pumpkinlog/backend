package engine

import (
	"github.com/pumpkinlog/backend/internal/domain"
)

type Engine struct {
	RuleEvaluator      RuleEvaluator
	ConditionEvaluator ConditionEvaluator
}

func NewEngine() *Engine {
	return &Engine{
		RuleEvaluator:      new(DefaultRuleEvaluator),
		ConditionEvaluator: new(DefaultConditionEvaluator),
	}
}

type EvaluateRegionsParams struct {
	Regions             []*domain.Region
	PresencesByRegion   map[string][]*domain.Presence
	RulesByRegion       map[string][]*domain.Rule
	ConditionsByRuleID  map[string][]*domain.Condition
	AnswerByConditionID map[string]*domain.Answer
}

func (e *Engine) EvaluateRegions(params *EvaluateRegionsParams) ([]*RegionEvaluation, error) {
	regionProfiles := make([]*RegionEvaluation, 0, len(params.Regions))

	for _, region := range params.Regions {
		rules := params.RulesByRegion[region.ID]
		ruleProfiles := make([]RuleEvaluation, 0, len(rules))

		for _, rule := range rules {
			conditions := params.ConditionsByRuleID[rule.ID]
			conditionProfiles := make([]ConditionEvaluation, 0, len(conditions))

			allConditionsPassed := true
			for _, condition := range conditions {
				answer := params.AnswerByConditionID[condition.ID]
				conditionProfile := e.ConditionEvaluator.Evaluate(condition, answer)
				conditionProfiles = append(conditionProfiles, conditionProfile)

				if !conditionProfile.Passed {
					allConditionsPassed = false
				}
			}

			ruleProfile := RuleEvaluation{
				Rule:       rule,
				Conditions: conditionProfiles,
				Passed:     false,
			}

			if allConditionsPassed {
				evalResult := e.RuleEvaluator.Evaluate(rule, params.PresencesByRegion[region.ID])
				ruleProfile.Conditions = conditionProfiles
				ruleProfile.Evaluation = evalResult
				ruleProfile.Passed = evalResult.Passed
			}

			ruleProfiles = append(ruleProfiles, ruleProfile)
		}

		regionProfiles = append(regionProfiles, &RegionEvaluation{
			Region: region,
			Rules:  ruleProfiles,
		})
	}

	return regionProfiles, nil
}

package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/pumpkinlog/backend/internal/app/engine"
	"github.com/pumpkinlog/backend/internal/app/period"
	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type EvaluationService struct {
	logger *slog.Logger
	engine *engine.Engine

	regionRepo        domain.RegionRepository
	ruleRepo          domain.RuleRepository
	ruleConditionRepo domain.RuleConditionRepository
	conditionRepo     domain.ConditionRepository
	answerRepo        domain.AnswerRepository
	evaluationRepo    domain.EvaluationRepository
	presenceRepo      domain.PresenceRepository
}

func NewEvaluationService(logger *slog.Logger, conn repository.Connection) domain.EvaluationService {
	return &EvaluationService{
		logger: logger,
		engine: engine.NewEngine(),

		regionRepo:        repository.NewPostgresRegionRepository(conn),
		ruleRepo:          repository.NewPostgresRuleRepository(conn),
		ruleConditionRepo: repository.NewPostgresRuleConditionRepository(conn),
		conditionRepo:     repository.NewPostgresConditionRepository(conn),
		answerRepo:        repository.NewPostgresAnswerRepository(conn),
		evaluationRepo:    repository.NewPostgresEvaluationRepository(conn),
		presenceRepo:      repository.NewPostgresPresenceRepository(conn),
	}
}

func (s *EvaluationService) prepareEvaluateRegionParams(ctx context.Context, userID, regionID string) (*engine.EvaluateRegionParams, error) {

	// 1. Load Region
	region, err := s.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("fetch region: %w", err)
	}

	// 2. Load Rules
	rules, err := s.ruleRepo.List(ctx, &domain.RuleFilter{
		RegionIDs: []string{regionID},
	})
	if err != nil {
		return nil, fmt.Errorf("fetch rules: %w", err)
	}

	// 3. Collect RuleIDs
	ruleIDs := make([]string, len(rules))
	for i, r := range rules {
		ruleIDs[i] = r.ID
	}

	// 4. Load RuleConditions using filter
	ruleConditionsSlice, err := s.ruleConditionRepo.List(ctx, &domain.RuleConditionFilter{
		RuleIDs: ruleIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch rule conditions: %w", err)
	}

	// 5. Organize by ruleID
	ruleConditions := make(map[string][]*domain.RuleCondition)
	conditionIDSet := make(map[string]struct{})

	for _, rc := range ruleConditionsSlice {
		ruleConditions[rc.RuleID] = append(ruleConditions[rc.RuleID], rc)
		conditionIDSet[rc.ConditionID] = struct{}{}
	}

	// 6. Extract unique condition IDs
	conditionIDs := make([]string, 0, len(conditionIDSet))
	for cid := range conditionIDSet {
		conditionIDs = append(conditionIDs, cid)
	}

	// 7. Load Conditions
	conditions, err := s.conditionRepo.List(ctx, &domain.ConditionFilter{
		ConditionIDs: conditionIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch conditions: %w", err)
	}

	// 8. Map Conditions by ID
	conditionMap := make(map[string]*domain.Condition, len(conditions))
	for _, cond := range conditions {
		conditionMap[cond.ID] = cond
	}

	// 9. Load Answers
	answers, err := s.answerRepo.List(ctx, userID, &domain.AnswerFilter{
		ConditionIDs: conditionIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch answers: %w", err)
	}

	// 10. Map Answers by ConditionID
	answerMap := make(map[string]*domain.Answer, len(answers))
	for _, a := range answers {
		answerMap[a.ConditionID] = a
	}

	// 11. Compute total period for the region
	prd := period.ComputeRulesPeriod(time.Now().UTC(), rules)

	// 12. Load Presences
	presences, err := s.presenceRepo.List(ctx, userID, &domain.PresenceFilter{
		RegionIDs: []string{regionID},
		Start:     &prd.Start,
		End:       &prd.End,
	})
	if err != nil {
		return nil, fmt.Errorf("evaluation service: list presences: %w", err)
	}

	// ✅ Final Params
	return &engine.EvaluateRegionParams{
		Region:         region,
		Presences:      presences,
		Rules:          rules,
		Conditions:     conditionMap,
		Answers:        answerMap,
		RuleConditions: ruleConditions,
	}, nil
}

func (s *EvaluationService) EvaluateRegion(ctx context.Context, userID, regionID string) (*domain.RegionEvaluation, error) {

	params, err := s.prepareEvaluateRegionParams(ctx, userID, regionID)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: prepare evaluate region params: %w", err)
	}

	evaluation, err := s.engine.EvaluateRegion(params)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: evaluate region: %w", err)
	}

	evaluation.UserID = userID
	evaluation.RegionID = regionID

	return evaluation, nil
}

func (s *EvaluationService) EvaluateRegions(ctx context.Context, userID string) ([]*domain.RegionEvaluation, error) {

	regions, err := s.regionRepo.List(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: list regions: %w", err)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	results := make([]*domain.RegionEvaluation, len(regions))

	for i, region := range regions {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int, region *domain.Region) {
			defer wg.Done()
			defer func() { <-sem }()

			evaluation, err := s.EvaluateRegion(ctx, userID, region.ID)
			if err != nil {
				s.logger.Error("evaluation service: evaluate region", "error", err)
				return
			}

			results[i] = evaluation
		}(i, region)
	}

	wg.Wait()
	close(sem)

	return results, nil
}

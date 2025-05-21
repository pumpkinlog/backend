package service

import (
	"context"
	"errors"
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

	regionRepo     domain.RegionRepository
	ruleRepo       domain.RuleRepository
	conditionRepo  domain.ConditionRepository
	answerRepo     domain.AnswerRepository
	evaluationRepo domain.EvaluationRepository
	presenceRepo   domain.PresenceRepository
}

func NewEvaluationService(logger *slog.Logger, conn repository.Connection) domain.EvaluationService {
	return &EvaluationService{
		logger: logger,
		engine: engine.NewEngine(),

		regionRepo:     repository.NewPostgresRegionRepository(conn),
		ruleRepo:       repository.NewPostgresRuleRepository(conn),
		conditionRepo:  repository.NewPostgresConditionRepository(conn),
		answerRepo:     repository.NewPostgresAnswerRepository(conn),
		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),
		presenceRepo:   repository.NewPostgresPresenceRepository(conn),
	}
}

func (s *EvaluationService) EvaluateRegion(ctx context.Context, userID, regionID string) (*domain.RegionEvaluation, error) {

	region, err := s.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: get region: %w", err)
	}

	evaluations, err := s.evaluateRegions(ctx, userID, []*domain.Region{region})
	if err != nil {
		return nil, fmt.Errorf("evaluation service: evaluate region: %w", err)
	}

	return evaluations[0], nil
}

func (s *EvaluationService) EvaluateRegions(ctx context.Context, userID string, regionIDs []string) ([]*domain.RegionEvaluation, error) {

	filter := &domain.RegionFilter{
		RegionIDs: regionIDs,
	}

	regions, err := s.regionRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: list regions: %w", err)
	}

	return s.evaluateRegions(ctx, userID, regions)
}

func (s *EvaluationService) EvaluateAllRegions(ctx context.Context, userID string) ([]*domain.RegionEvaluation, error) {
	return s.evaluateRegions(ctx, userID, nil)
}

func (s *EvaluationService) evaluateRegions(ctx context.Context, userID string, regions []*domain.Region) ([]*domain.RegionEvaluation, error) {

	cached := make([]*domain.RegionEvaluation, 0, len(regions))
	params := &engine.EvaluateRegionsParams{
		Regions:             make([]*domain.Region, 0),
		PresencesByRegion:   make(map[string][]*domain.Presence),
		RulesByRegion:       make(map[string][]*domain.Rule),
		ConditionsByRuleID:  make(map[string][]*domain.Condition),
		AnswerByConditionID: make(map[string]*domain.Answer),
	}

	for _, region := range regions {
		evaluation, err := s.evaluationRepo.GetByID(ctx, userID, region.ID)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("evaluation service: fetch evaluation: %w", err)
		}

		if evaluation != nil {
			s.logger.Info("evaluation found, skipping evaluation", "userID", userID, "regionID", region.ID)
			cached = append(cached, evaluation)
			continue
		}

		s.logger.Info("no evaluation found, creating a new one", "userID", userID, "regionID", region.ID)
		params.Regions = append(params.Regions, region)
	}

	ruleFilter := &domain.RuleFilter{
		RegionIDs: make([]string, len(params.Regions)),
	}

	for i, region := range params.Regions {
		ruleFilter.RegionIDs[i] = region.ID
	}

	rules, err := s.ruleRepo.List(ctx, ruleFilter)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: list rules: %w", err)
	}

	conditionsFilter := &domain.ConditionFilter{
		RuleIDs: make([]string, 0),
	}

	for _, rule := range rules {
		conditionsFilter.RuleIDs = append(conditionsFilter.RuleIDs, rule.ID)
		params.RulesByRegion[rule.RegionID] = append(params.RulesByRegion[rule.RegionID], rule)
	}

	conditions, err := s.conditionRepo.List(ctx, conditionsFilter)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: list conditions: %w", err)
	}

	answersFilter := &domain.AnswerFilter{
		ConditionIDs: make([]string, 0),
	}

	for _, condition := range conditions {
		answersFilter.ConditionIDs = append(answersFilter.ConditionIDs, condition.ID)
		params.ConditionsByRuleID[condition.RuleID] = append(params.ConditionsByRuleID[condition.RuleID], condition)

	}

	answers, err := s.answerRepo.List(ctx, userID, answersFilter)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: list answers: %w", err)
	}

	for _, answer := range answers {
		params.AnswerByConditionID[answer.ConditionID] = answer
	}

	prd := period.ComputePeriodByRegion(time.Now().UTC(), rules)

	presences, err := s.presenceRepo.ListByRegionBounds(ctx, userID, prd)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: list presences: %w", err)
	}

	for _, presence := range presences {
		params.PresencesByRegion[presence.RegionID] = append(params.PresencesByRegion[presence.RegionID], presence)
	}

	re, err := s.engine.EvaluateRegions(params)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: evaluate: %w", err)
	}

	evaluations := make([]*domain.RegionEvaluation, len(re))
	evaluatedAt := time.Now().UTC()

	for i, e := range re {
		results, passed := convertToRuleEvaluations(e.Rules)

		evaluations[i] = &domain.RegionEvaluation{
			UserID:      userID,
			RegionID:    e.Region.ID,
			Passed:      passed,
			Evaluations: results,
			EvaluatedAt: evaluatedAt,
		}
	}

	var wg sync.WaitGroup

	for _, evaluation := range evaluations {
		wg.Add(1)

		go func(evaluation *domain.RegionEvaluation) {
			defer wg.Done()

			if err := s.evaluationRepo.CreateOrUpdate(ctx, evaluation); err != nil {
				fmt.Printf("evaluation service: create or update evaluation: %v\n", err)
			}
		}(evaluation)
	}

	wg.Wait()

	for _, evaluation := range evaluations {
		cached = append(cached, evaluation)
	}

	return cached, nil
}

func convertToRuleEvaluations(rules []engine.RuleEvaluation) ([]domain.RuleEvaluation, bool) {
	evaluations := make([]domain.RuleEvaluation, 0, len(rules))
	var passed bool

	for _, rule := range rules {
		evaluations = append(evaluations, domain.RuleEvaluation{
			Passed: rule.Passed,
			Rule:   rule.Rule,
			Logic: domain.RuleLogicEvaluation{
				Resident:       rule.Evaluation.Passed,
				Count:          rule.Evaluation.Count,
				Remaining:      rule.Evaluation.Remaining,
				Start:          rule.Evaluation.Start,
				End:            rule.Evaluation.End,
				ConsecutiveEnd: rule.Evaluation.ConsecutiveEnd,
				Metadata:       rule.Evaluation.Metadata,
			},
			Conditions: convertToConditionEvaluations(rule.Conditions),
		})
	}

	return evaluations, passed
}

func convertToConditionEvaluations(conditions []engine.ConditionEvaluation) []domain.ConditionEvaluation {
	evaluations := make([]domain.ConditionEvaluation, 0, len(conditions))

	for _, condition := range conditions {
		evaluations = append(evaluations, domain.ConditionEvaluation{
			Condition: condition.Condition,
			Answer:    condition.Answer,
			Passed:    condition.Passed,
			Skipped:   condition.Skipped,
		})
	}

	return evaluations
}

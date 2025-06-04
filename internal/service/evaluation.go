package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pumpkinlog/backend/internal/app/engine"
	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type EvaluationService struct {
	logger *slog.Logger
	engine *engine.Engine

	regionRepo        domain.RegionRepository
	ruleRepo          domain.RuleRepository
	conditionRepo     domain.ConditionRepository
	ruleConditionRepo domain.RuleConditionRepository
	answerRepo        domain.AnswerRepository
	evaluationRepo    domain.EvaluationRepository
	presenceRepo      domain.PresenceRepository
}

func NewEvaluationService(logger *slog.Logger, conn *pgxpool.Pool) domain.EvaluationService {
	return &EvaluationService{
		logger: logger,
		engine: engine.NewEngine(),

		regionRepo:        repository.NewPostgresRegionRepository(conn),
		ruleRepo:          repository.NewPostgresRuleRepository(conn),
		conditionRepo:     repository.NewPostgresConditionRepository(conn),
		ruleConditionRepo: repository.NewPostgresRuleConditionRepository(conn),
		answerRepo:        repository.NewPostgresAnswerRepository(conn),
		evaluationRepo:    repository.NewPostgresEvaluationRepository(conn),
		presenceRepo:      repository.NewPostgresPresenceRepository(conn),
	}
}

func (s *EvaluationService) EvaluateRegion(ctx context.Context, userID int64, regionID string) (*domain.RegionEvaluation, error) {

	evaluation, err := s.evaluationRepo.GetByID(ctx, userID, regionID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("get evaluation by ID: %w", err)
	}

	if evaluation != nil {
		// Evaluation already computed
		return evaluation, nil
	}

	s.logger.Debug("running region evaluation", "userId", userID, "regionId", regionID)

	aggregate, err := s.loadAggregate(ctx, userID, regionID)
	if err != nil {
		return nil, fmt.Errorf("load aggregate: %w", err)
	}

	evaluations, passed, err := s.engine.EvaluateRegion(aggregate)
	if err != nil {
		return nil, fmt.Errorf("evaluation service: evaluate region: %w", err)
	}

	evaluation = &domain.RegionEvaluation{
		UserID:          userID,
		RegionID:        regionID,
		Passed:          passed,
		RuleEvaluations: evaluations,
		EvaluatedAt:     time.Now().UTC(),
	}

	if err := s.evaluationRepo.CreateOrUpdate(ctx, evaluation); err != nil {
		return nil, fmt.Errorf("create or update evaluation: %w", err)
	}

	return evaluation, nil
}

func (s *EvaluationService) loadAggregate(ctx context.Context, userID int64, regionID string) (*domain.RegionAggregate, error) {

	region, err := s.regionRepo.GetByID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("get region: %w", err)
	}

	rules, err := s.ruleRepo.ListByRegionID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("list rules by region ID: %w", err)
	}

	conditions, err := s.conditionRepo.ListByRegionID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("list conditions by region ID: %w", err)
	}

	ruleConditions, err := s.ruleConditionRepo.ListByRegionID(ctx, region.ID)
	if err != nil {
		return nil, fmt.Errorf("list rule conditions by region ID: %w", err)
	}

	answers, err := s.answerRepo.ListByRegionID(ctx, userID, region.ID)
	if err != nil {
		return nil, fmt.Errorf("list answers by region ID: %w", err)
	}

	start, end, err := s.maxPeriodForRules(time.Now().UTC(), rules)
	if err != nil {
		return nil, fmt.Errorf("get max period for rules: %w", err)
	}

	presences, err := s.presenceRepo.ListByRegionPeriod(ctx, userID, region.ID, start, end)
	if err != nil {
		return nil, fmt.Errorf("list presences by region period: %w", err)
	}

	aggregate := &domain.RegionAggregate{
		Region:         region,
		Presences:      presences,
		Rules:          rules,
		RuleConditions: ruleConditions,
		Conditions:     sliceToMap(conditions, func(c *domain.Condition) int64 { return c.ID }),
		Answers:        sliceToMap(answers, func(a *domain.Answer) int64 { return a.ConditionID }),
	}

	return aggregate, nil
}

func (*EvaluationService) maxPeriodForRules(asOf time.Time, rules []*domain.Rule) (time.Time, time.Time, error) {
	var start, end time.Time

	for i, rule := range rules {
		s, f, err := rule.Period(asOf)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("get max period for rule %d: %w", rule.ID, err)
		}

		if i == 0 {
			start, end = s, f
			continue
		}

		if s.Before(start) {
			start = s
		}

		if f.After(end) {
			end = f
		}
	}

	if start.IsZero() || end.IsZero() {
		return time.Time{}, time.Time{}, errors.New("no valid period found for rules")
	}

	return start, end, nil
}

func sliceToMap[T any, K comparable](items []*T, keyFn func(*T) K) map[K]*T {
	m := make(map[K]*T, len(items))
	for _, item := range items {
		m[keyFn(item)] = item
	}
	return m
}

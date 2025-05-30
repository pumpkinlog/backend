package seed

import (
	"context"
	"log/slog"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type Seeder struct {
	logger *slog.Logger

	regionRepo          domain.RegionRepository
	ruleRepo            domain.RuleRepository
	conditionRepo       domain.ConditionRepository
	regionConditionRepo domain.RegionConditionRepository
	ruleConditionRepo   domain.RuleConditionRepository
}

func NewSeeder(logger *slog.Logger, conn repository.Connection) *Seeder {
	return &Seeder{
		logger: logger,

		regionRepo:          repository.NewPostgresRegionRepository(conn),
		ruleRepo:            repository.NewPostgresRuleRepository(conn),
		conditionRepo:       repository.NewPostgresConditionRepository(conn),
		regionConditionRepo: repository.NewPostgresRegionConditionRepository(conn),
		ruleConditionRepo:   repository.NewPostgresRuleConditionRepository(conn),
	}
}

func (s *Seeder) Seed(ctx context.Context) {

	start := time.Now()

	for _, region := range s.regions() {

		if region.YearStartMonth == 0 {
			region.YearStartMonth = 1
		}

		if region.YearStartDay == 0 {
			region.YearStartDay = 1
		}

		if err := s.regionRepo.CreateOrUpdate(ctx, region); err != nil {
			s.logger.Error("failed to upsert region", "region_id", region.ID, "error", err)
		}
	}

	for _, rule := range s.rules() {
		if err := s.ruleRepo.CreateOrUpdate(ctx, rule); err != nil {
			s.logger.Error("failed to upsert rule", "rule_id", rule.ID, "error", err)
		}
	}

	for _, cond := range s.conditions() {
		if err := s.conditionRepo.CreateOrUpdate(ctx, cond); err != nil {
			s.logger.Error("failed to upsert condition", "condition_id", cond.ID, "error", err)
		}
	}

	for _, rc := range s.regionConditions() {
		if err := s.regionConditionRepo.CreateOrUpdate(ctx, rc); err != nil {
			s.logger.Error("failed to upsert region condition", "region_id", rc.RegionID, "condition_id", rc.ConditionID, "error", err)
		}
	}

	for _, rc := range s.ruleConditions() {
		if err := s.ruleConditionRepo.CreateOrUpdate(ctx, rc); err != nil {
			s.logger.Error("failed to upsert rule condition", "rule_id", rc.RuleID, "condition_id", rc.ConditionID, "error", err)
		}
	}

	s.logger.Info("seeding completed", "duration", time.Since(start).String())
}

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

	regionRepo    domain.RegionRepository
	ruleRepo      domain.RuleRepository
	conditionRepo domain.ConditionRepository
}

func NewSeeder(logger *slog.Logger, conn repository.Connection) *Seeder {
	return &Seeder{
		logger: logger,

		regionRepo:    repository.NewPostgresRegionRepository(conn),
		ruleRepo:      repository.NewPostgresRuleRepository(conn),
		conditionRepo: repository.NewPostgresConditionRepository(conn),
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

	for _, condition := range s.conditions() {
		if err := s.conditionRepo.CreateOrUpdate(ctx, condition); err != nil {
			s.logger.Error("failed to upsert condition", "condition_id", condition.ID, "error", err)
		}
	}

	s.logger.Info("seeding completed", "duration", time.Since(start).String())
}

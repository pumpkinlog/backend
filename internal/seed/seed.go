package seed

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/service"
)

type Seeder struct {
	logger *slog.Logger

	regionSvc        domain.RegionService
	ruleSvc          domain.RuleService
	conditionSvc     domain.ConditionService
	ruleConditionSvc domain.RuleConditionService
}

func NewSeeder(logger *slog.Logger, conn *pgxpool.Pool) *Seeder {
	return &Seeder{
		logger: logger,

		regionSvc:        service.NewRegionService(logger, conn),
		ruleSvc:          service.NewRuleService(logger, conn),
		conditionSvc:     service.NewConditionService(logger, conn),
		ruleConditionSvc: service.NewRuleConditionService(logger, conn),
	}
}

func (s *Seeder) Seed(ctx context.Context) {

	start := time.Now()

	for _, region := range regions {

		if region.YearStartMonth == 0 {
			region.YearStartMonth = 1
		}

		if region.YearStartDay == 0 {
			region.YearStartDay = 1
		}

		if err := s.regionSvc.CreateOrUpdate(ctx, &region); err != nil {
			s.logger.Error("failed to upsert region", "region_id", region.ID, "error", err)
		}
	}

	for _, rule := range rules {
		if err := s.ruleSvc.CreateOrUpdate(ctx, &rule); err != nil {
			s.logger.Error("failed to upsert rule", "rule_id", rule.ID, "error", err)
		}
	}

	for _, cond := range conditions {
		if err := s.conditionSvc.CreateOrUpdate(ctx, &cond); err != nil {
			s.logger.Error("failed to upsert condition", "condition_id", cond.ID, "error", err)
		}
	}

	for _, rc := range ruleConditions {
		if err := s.ruleConditionSvc.CreateOrUpdate(ctx, &rc); err != nil {
			s.logger.Error("failed to upsert rule condition", "rule_id", rc.RuleID, "condition_id", rc.ConditionID, "error", err)
		}
	}

	s.logger.Info("seeding completed", "duration", time.Since(start).String())
}

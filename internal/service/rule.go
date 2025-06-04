package service

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type RuleService struct {
	logger *slog.Logger

	ruleRepo       domain.RuleRepository
	evaluationRepo domain.EvaluationRepository
}

func NewRuleService(logger *slog.Logger, conn *pgxpool.Pool) domain.RuleService {
	return &RuleService{
		logger: logger,

		ruleRepo:       repository.NewPostgresRuleRepository(conn),
		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),
	}
}

func (s *RuleService) CreateOrUpdate(ctx context.Context, rule *domain.Rule) error {

	if err := rule.Validate(); err != nil {
		return err
	}

	// @TODO: Tx?

	if err := s.ruleRepo.CreateOrUpdate(ctx, rule); err != nil {
		return err
	}

	// Delete existing evaluations for the region
	if err := s.evaluationRepo.DeleteByRegionID(ctx, rule.RegionID); err != nil {
		return err
	}

	s.logger.Debug("cleared stale evaluations", "regionId", rule.RegionID, "ruleId", rule.ID)

	return nil
}

package service

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type RuleConditionService struct {
	logger *slog.Logger

	ruleConditionRepo domain.RuleConditionRepository
	evaluationRepo    domain.EvaluationRepository
}

func NewRuleConditionService(logger *slog.Logger, conn *pgxpool.Pool) domain.RuleConditionService {
	return &RuleConditionService{
		logger: logger,

		ruleConditionRepo: repository.NewPostgresRuleConditionRepository(conn),
		evaluationRepo:    repository.NewPostgresEvaluationRepository(conn),
	}
}

func (s *RuleConditionService) CreateOrUpdate(ctx context.Context, ruleCondition *domain.RuleCondition) error {

	if err := ruleCondition.Validate(); err != nil {
		return err
	}

	if err := s.ruleConditionRepo.CreateOrUpdate(ctx, ruleCondition); err != nil {
		return err
	}

	// Delete existing evaluations for the region
	if err := s.evaluationRepo.DeleteByRegionID(ctx, ruleCondition.RegionID); err != nil {
		return err
	}

	s.logger.Debug("cleared stale evaluations", "regionId", ruleCondition.RegionID, "ruleId", ruleCondition.RuleID, "conditionId", ruleCondition.ConditionID)

	return nil
}

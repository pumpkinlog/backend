package service

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type ConditionService struct {
	logger *slog.Logger

	conditionRepo  domain.ConditionRepository
	evaluationRepo domain.EvaluationRepository
}

func NewConditionService(logger *slog.Logger, conn *pgxpool.Pool) domain.ConditionService {
	return &ConditionService{
		logger: logger,

		conditionRepo:  repository.NewPostgresConditionRepository(conn),
		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),
	}
}

func (s *ConditionService) CreateOrUpdate(ctx context.Context, condition *domain.Condition) error {

	if err := condition.Validate(); err != nil {
		return err
	}

	if err := s.conditionRepo.CreateOrUpdate(ctx, condition); err != nil {
		return err
	}

	// Delete existing evaluations for the region
	if err := s.evaluationRepo.DeleteByRegionID(ctx, condition.RegionID); err != nil {
		return err
	}

	s.logger.Debug("cleared stale evaluations", "regionId", condition.RegionID, "conditionId", condition.ID)

	return nil
}

package service

import (
	"context"
	"log/slog"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type ConditionService struct {
	logger *slog.Logger

	conditionRepo  domain.ConditionRepository
	evaluationRepo domain.EvaluationRepository
}

func NewConditionService(logger *slog.Logger, conn repository.Connection) domain.ConditionService {
	return &ConditionService{
		logger: logger,

		conditionRepo:  repository.NewPostgresConditionRepository(conn),
		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),
	}
}

func (s ConditionService) GetByID(ctx context.Context, conditionID domain.Code) (*domain.Condition, error) {
	if err := conditionID.Validate(); err != nil {
		return nil, err
	}

	return s.conditionRepo.GetByID(ctx, conditionID)
}

func (s ConditionService) List(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error) {
	if filter == nil {
		filter = &domain.ConditionFilter{}
	}

	return s.conditionRepo.List(ctx, filter)
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

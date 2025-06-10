package service

import (
	"context"
	"log/slog"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type RuleService struct {
	logger *slog.Logger

	ruleRepo       domain.RuleRepository
	evaluationRepo domain.EvaluationRepository
}

func NewRuleService(logger *slog.Logger, conn repository.Connection) domain.RuleService {
	return &RuleService{
		logger: logger,

		ruleRepo:       repository.NewPostgresRuleRepository(conn),
		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),
	}
}

func (s *RuleService) GetByID(ctx context.Context, ruleID domain.Code) (*domain.Rule, error) {
	if err := ruleID.Validate(); err != nil {
		return nil, err
	}

	return s.ruleRepo.GetByID(ctx, ruleID)
}

func (s *RuleService) List(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error) {
	if filter == nil {
		filter = &domain.RuleFilter{}
	}

	return s.ruleRepo.List(ctx, filter)
}

func (s *RuleService) CreateOrUpdate(ctx context.Context, rule *domain.Rule) error {
	if err := rule.Validate(); err != nil {
		return err
	}

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

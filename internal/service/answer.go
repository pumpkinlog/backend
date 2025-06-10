package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type AnswerService struct {
	logger *slog.Logger

	conditionRepo  domain.ConditionRepository
	answerRepo     domain.AnswerRepository
	evaluationRepo domain.EvaluationRepository
}

func NewAnswerService(logger *slog.Logger, conn repository.Connection) domain.AnswerService {
	return &AnswerService{
		logger: logger,

		conditionRepo:  repository.NewPostgresConditionRepository(conn),
		answerRepo:     repository.NewPostgresAnswerRepository(conn),
		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),
	}
}

func (s *AnswerService) GetByID(ctx context.Context, userID int64, conditionID domain.Code) (*domain.Answer, error) {
	if userID <= 0 {
		return nil, domain.ValidationError("user ID is required")
	}

	if err := conditionID.Validate(); err != nil {
		return nil, err
	}

	return s.answerRepo.GetByID(ctx, userID, conditionID)
}

func (s *AnswerService) CreateOrUpdate(ctx context.Context, userID int64, conditionID domain.Code, value any) error {
	now := time.Now().UTC()

	condition, err := s.conditionRepo.GetByID(ctx, conditionID)
	if err != nil {
		return err
	}

	answer := &domain.Answer{
		UserID:      userID,
		ConditionID: conditionID,
		RegionID:    condition.RegionID,
		Value:       value,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := answer.Validate(); err != nil {
		return err
	}

	if err := s.answerRepo.CreateOrUpdate(ctx, answer); err != nil {
		return fmt.Errorf("create or update answer: %w", err)
	}

	// Delete existing evaluations for the user region
	if err := s.evaluationRepo.DeleteByUserAndRegionID(ctx, userID, condition.RegionID); err != nil {
		return err
	}

	s.logger.Debug("cleared stale evaluations", "userId", userID, "regionId", condition.RegionID, "conditionId", conditionID)

	return nil
}

func (s *AnswerService) Delete(ctx context.Context, userID int64, conditionID domain.Code) error {
	if userID <= 0 {
		return domain.ValidationError("user ID is required")
	}

	if err := conditionID.Validate(); err != nil {
		return err
	}

	answer, err := s.answerRepo.GetByID(ctx, userID, conditionID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return fmt.Errorf("get answer by ID: %w", err)
	}

	if err := s.answerRepo.Delete(ctx, userID, conditionID); err != nil {
		return fmt.Errorf("delete answer: %w", err)
	}

	// Delete existing evaluations for the user region
	if err := s.evaluationRepo.DeleteByUserAndRegionID(ctx, userID, answer.RegionID); err != nil {
		return err
	}

	s.logger.Debug("cleared stale evaluations", "userId", userID, "regionId", answer.RegionID, "conditionId", conditionID)

	return nil
}

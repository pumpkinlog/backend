package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type AnswerService struct {
	conditionRepo domain.ConditionRepository
	answerRepo    domain.AnswerRepository
}

func NewAnswerService(conn *pgxpool.Pool) domain.AnswerService {
	return &AnswerService{
		answerRepo: repository.NewPostgresAnswerRepository(conn),
	}
}

func (s *AnswerService) CreateOrUpdate(ctx context.Context, userID, conditionID int64, value any) error {
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

	return s.answerRepo.CreateOrUpdate(ctx, answer)
}

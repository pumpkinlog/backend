package service

import (
	"context"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type AnswerService struct {
	answerRepo domain.AnswerRepository
}

func NewAnswerService(conn repository.Connection) domain.AnswerService {
	return &AnswerService{
		answerRepo: repository.NewPostgresAnswerRepository(conn),
	}
}

func (s *AnswerService) CreateOrUpdate(ctx context.Context, userID, conditionID string, value any) error {
	now := time.Now().UTC()

	answer := &domain.Answer{
		UserID:      userID,
		ConditionID: conditionID,
		Value:       value,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := answer.Validate(); err != nil {
		return err
	}

	return s.answerRepo.CreateOrUpdate(ctx, answer)
}

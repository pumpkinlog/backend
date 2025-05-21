package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type AnswerRepo struct {
	GetByIDFunc        func(ctx context.Context, userID, conditionID string) (*domain.Answer, error)
	ListFunc           func(ctx context.Context, userID string, filter *domain.AnswerFilter) ([]*domain.Answer, error)
	CreateOrUpdateFunc func(ctx context.Context, answer *domain.Answer) error
}

func (m AnswerRepo) GetByID(ctx context.Context, userID, conditionID string) (*domain.Answer, error) {
	return m.GetByIDFunc(ctx, userID, conditionID)
}

func (m AnswerRepo) List(ctx context.Context, userID string, filter *domain.AnswerFilter) ([]*domain.Answer, error) {
	return m.ListFunc(ctx, userID, filter)
}

func (m AnswerRepo) CreateOrUpdate(ctx context.Context, answer *domain.Answer) error {
	return m.CreateOrUpdateFunc(ctx, answer)
}

package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type AnswerRepository struct {
	GetByIDFunc        func(ctx context.Context, userID, conditionID string) (*domain.Answer, error)
	ListFunc           func(ctx context.Context, userID string, filter *domain.AnswerFilter) ([]*domain.Answer, error)
	CreateOrUpdateFunc func(ctx context.Context, answer *domain.Answer) error
	DeleteFunc         func(ctx context.Context, userID, conditionID string) error
}

func (m AnswerRepository) GetByID(ctx context.Context, userID, conditionID string) (*domain.Answer, error) {
	return m.GetByIDFunc(ctx, userID, conditionID)
}

func (m AnswerRepository) List(ctx context.Context, userID string, filter *domain.AnswerFilter) ([]*domain.Answer, error) {
	return m.ListFunc(ctx, userID, filter)
}

func (m AnswerRepository) CreateOrUpdate(ctx context.Context, answer *domain.Answer) error {
	return m.CreateOrUpdateFunc(ctx, answer)
}

func (m AnswerRepository) Delete(ctx context.Context, userID, conditionID string) error {
	return m.DeleteFunc(ctx, userID, conditionID)
}

type AnswerService struct {
	CreateOrUpdateFunc func(ctx context.Context, userID, conditionID string, value any) error
}

func (m AnswerService) CreateOrUpdate(ctx context.Context, userID, conditionID string, value any) error {
	return m.CreateOrUpdateFunc(ctx, userID, conditionID, value)
}

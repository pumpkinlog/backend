package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type AnswerRepository struct {
	GetByIDFunc        func(ctx context.Context, userID, conditionID int64) (*domain.Answer, error)
	ListByRegionIDFunc func(ctx context.Context, userID int64, regionID string) ([]*domain.Answer, error)
	CreateOrUpdateFunc func(ctx context.Context, answer *domain.Answer) error
	DeleteFunc         func(ctx context.Context, userID, conditionID int64) error
}

func (m AnswerRepository) GetByID(ctx context.Context, userID, conditionID int64) (*domain.Answer, error) {
	return m.GetByIDFunc(ctx, userID, conditionID)
}

func (m AnswerRepository) ListByRegionID(ctx context.Context, userID int64, regionID string) ([]*domain.Answer, error) {
	return m.ListByRegionIDFunc(ctx, userID, regionID)
}

func (m AnswerRepository) CreateOrUpdate(ctx context.Context, answer *domain.Answer) error {
	return m.CreateOrUpdateFunc(ctx, answer)
}

func (m AnswerRepository) Delete(ctx context.Context, userID, conditionID int64) error {
	return m.DeleteFunc(ctx, userID, conditionID)
}

type AnswerService struct {
	CreateOrUpdateFunc func(ctx context.Context, userID, conditionID int64, value any) error
}

func (m AnswerService) CreateOrUpdate(ctx context.Context, userID, conditionID int64, value any) error {
	return m.CreateOrUpdateFunc(ctx, userID, conditionID, value)
}

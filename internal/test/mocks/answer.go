package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type AnswerRepository struct {
	GetByIDFunc        func(ctx context.Context, userID int64, conditionID domain.Code) (*domain.Answer, error)
	ListByRegionIDFunc func(ctx context.Context, userID int64, regionID domain.RegionID) ([]*domain.Answer, error)
	CreateOrUpdateFunc func(ctx context.Context, answer *domain.Answer) error
	DeleteFunc         func(ctx context.Context, userID int64, conditionID domain.Code) error
}

func (m AnswerRepository) GetByID(ctx context.Context, userID int64, conditionID domain.Code) (*domain.Answer, error) {
	return m.GetByIDFunc(ctx, userID, conditionID)
}

func (m AnswerRepository) ListByRegionID(ctx context.Context, userID int64, regionID domain.RegionID) ([]*domain.Answer, error) {
	return m.ListByRegionIDFunc(ctx, userID, regionID)
}

func (m AnswerRepository) CreateOrUpdate(ctx context.Context, answer *domain.Answer) error {
	return m.CreateOrUpdateFunc(ctx, answer)
}

func (m AnswerRepository) Delete(ctx context.Context, userID int64, conditionID domain.Code) error {
	return m.DeleteFunc(ctx, userID, conditionID)
}

type AnswerService struct {
	GetByIDFunc        func(ctx context.Context, userID int64, conditionID domain.Code) (*domain.Answer, error)
	CreateOrUpdateFunc func(ctx context.Context, userID int64, conditionID domain.Code, value any) error
	DeleteFunc         func(ctx context.Context, userID int64, conditionID domain.Code) error
}

func (m AnswerService) GetByID(ctx context.Context, userID int64, conditionID domain.Code) (*domain.Answer, error) {
	return m.GetByIDFunc(ctx, userID, conditionID)
}

func (m AnswerService) CreateOrUpdate(ctx context.Context, userID int64, conditionID domain.Code, value any) error {
	return m.CreateOrUpdateFunc(ctx, userID, conditionID, value)
}

func (m AnswerService) Delete(ctx context.Context, userID int64, conditionID domain.Code) error {
	return m.DeleteFunc(ctx, userID, conditionID)
}

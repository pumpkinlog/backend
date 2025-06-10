package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type UserRepository struct {
	GetByIDFunc func(ctx context.Context, userID int64) (*domain.User, error)
	CreateFunc  func(ctx context.Context, user *domain.User) error
	UpdateFunc  func(ctx context.Context, user *domain.User) error
	DeleteFunc  func(ctx context.Context, userID int64) error
}

func (m UserRepository) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	return m.GetByIDFunc(ctx, userID)
}

func (m UserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.CreateFunc(ctx, user)
}

func (m UserRepository) Update(ctx context.Context, user *domain.User) error {
	return m.UpdateFunc(ctx, user)
}

func (m UserRepository) Delete(ctx context.Context, userID int64) error {
	return m.DeleteFunc(ctx, userID)
}

type UserService struct {
	GetByIDFunc func(ctx context.Context, userID int64) (*domain.User, error)
	CreateFunc  func(ctx context.Context, favoriteRegions, wantResidency []domain.RegionID) error
	UpdateFunc  func(ctx context.Context, userID int64, favoriteRegions, wantResidency []domain.RegionID) error
}

func (m UserService) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	return m.GetByIDFunc(ctx, userID)
}

func (m UserService) Create(ctx context.Context, favoriteRegions, wantResidency []domain.RegionID) error {
	return m.CreateFunc(ctx, favoriteRegions, wantResidency)
}

func (m UserService) Update(ctx context.Context, userID int64, favorioteRegions, wantResidency []domain.RegionID) error {
	return m.UpdateFunc(ctx, userID, favorioteRegions, wantResidency)
}

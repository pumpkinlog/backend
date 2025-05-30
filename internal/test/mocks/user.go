package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type UserRepository struct {
	GetByIDFunc func(ctx context.Context, userID string) (*domain.User, error)
	CreateFunc  func(ctx context.Context, user *domain.User) error
	UpdateFunc  func(ctx context.Context, user *domain.User) error
	DeleteFunc  func(ctx context.Context, userID string) error
}

func (m UserRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	return m.GetByIDFunc(ctx, userID)
}

func (m UserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.CreateFunc(ctx, user)
}

func (m UserRepository) Update(ctx context.Context, user *domain.User) error {
	return m.UpdateFunc(ctx, user)
}

func (m UserRepository) Delete(ctx context.Context, userID string) error {
	return m.DeleteFunc(ctx, userID)
}

type UserService struct {
	CreateFunc func(ctx context.Context, userID string) error
	UpdateFunc func(ctx context.Context, userID string, favoriteRegions, wantResidency []string) error
}

func (m UserService) Create(ctx context.Context, userID string) error {
	return m.CreateFunc(ctx, userID)
}

func (m UserService) Update(ctx context.Context, userID string, favorioteRegions, wantResidency []string) error {
	return m.UpdateFunc(ctx, userID, favorioteRegions, wantResidency)
}

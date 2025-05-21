package mocks

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type UserRepo struct {
	GetByIDFunc func(ctx context.Context, id string) (*domain.User, error)
	CreateFunc  func(ctx context.Context, user *domain.User) error
	UpdateFunc  func(ctx context.Context, user *domain.User) error
	DeleteFunc  func(ctx context.Context, id string) error
}

func (m UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m UserRepo) Create(ctx context.Context, user *domain.User) error {
	return m.CreateFunc(ctx, user)
}

func (m UserRepo) Update(ctx context.Context, user *domain.User) error {
	return m.UpdateFunc(ctx, user)
}

func (m UserRepo) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

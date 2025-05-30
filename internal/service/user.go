package service

import (
	"context"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(conn repository.Connection) domain.UserService {
	return &UserService{
		userRepo: repository.NewPostgresUserRepository(conn),
	}
}

func (s *UserService) Create(ctx context.Context, userID string) error {
	now := time.Now().UTC()

	user := &domain.User{
		ID:              userID,
		FavoriteRegions: make([]string, 0),
		WantResidency:   make([]string, 0),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := user.Validate(); err != nil {
		return err
	}

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, userID string, favoriteRegions, wantResidency []string) error {

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if len(favoriteRegions) > 0 {
		user.FavoriteRegions = favoriteRegions
	}

	if len(wantResidency) > 0 {
		user.WantResidency = wantResidency
	}

	user.UpdatedAt = time.Now().UTC()

	if err := user.Validate(); err != nil {
		return err
	}

	return s.userRepo.Update(ctx, user)
}

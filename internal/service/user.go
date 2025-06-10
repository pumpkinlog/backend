package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type UserService struct {
	logger *slog.Logger

	userRepo domain.UserRepository
}

func NewUserService(logger *slog.Logger, conn repository.Connection) domain.UserService {
	return &UserService{
		logger: logger,

		userRepo: repository.NewPostgresUserRepository(conn),
	}
}

func (s *UserService) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	if userID < 0 {
		return nil, fmt.Errorf("%w: user ID cannot be negative", domain.ErrValidation)
	}

	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) Create(ctx context.Context, favoriteRegions, wantResidency []domain.RegionID) error {
	now := time.Now().UTC()

	user := &domain.User{
		FavoriteRegions: favoriteRegions,
		WantResidency:   wantResidency,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if user.FavoriteRegions == nil {
		user.FavoriteRegions = make([]domain.RegionID, 0)
	}

	if user.WantResidency == nil {
		user.WantResidency = make([]domain.RegionID, 0)
	}

	if err := user.Validate(); err != nil {
		return err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	s.logger.Debug("created new user", "userId", user.ID)

	return nil
}

func (s *UserService) Update(ctx context.Context, userID int64, favoriteRegions, wantResidency []domain.RegionID) error {

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

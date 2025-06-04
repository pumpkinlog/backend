package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
)

type UserService struct {
	logger *slog.Logger

	userRepo domain.UserRepository
}

func NewUserService(logger *slog.Logger, conn *pgxpool.Pool) domain.UserService {
	return &UserService{
		logger: logger,

		userRepo: repository.NewPostgresUserRepository(conn),
	}
}

func (s *UserService) Create(ctx context.Context, favoriteRegions, wantResidency []string) error {
	now := time.Now().UTC()

	user := &domain.User{
		FavoriteRegions: favoriteRegions,
		WantResidency:   wantResidency,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if user.FavoriteRegions == nil {
		user.FavoriteRegions = make([]string, 0)
	}

	if user.WantResidency == nil {
		user.WantResidency = make([]string, 0)
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

func (s *UserService) Update(ctx context.Context, userID int64, favoriteRegions, wantResidency []string) error {

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

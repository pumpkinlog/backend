package domain

import (
	"context"
	"fmt"
	"time"
)

var maxRegions = 300

type User struct {
	ID              int64      `json:"id"`
	FavoriteRegions []RegionID `json:"favoriteRegions"`
	WantResidency   []RegionID `json:"wantResidency"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

func (u *User) Validate() error {

	if u.ID < 0 {
		return fmt.Errorf("%w: id is invalid", ErrValidation)
	}

	if len(u.FavoriteRegions) > maxRegions {
		return fmt.Errorf("%w: favorite regions cannot be greater than %d", ErrValidation, maxRegions)
	}

	if len(u.WantResidency) > maxRegions {
		return fmt.Errorf("%w: want residency cannot be greater than %d", ErrValidation, maxRegions)
	}

	if u.CreatedAt.IsZero() {
		return fmt.Errorf("%w: created at timestamp is invalid", ErrValidation)
	}

	if u.UpdatedAt.IsZero() {
		return fmt.Errorf("%w: updated at timestamp is invalid", ErrValidation)
	}

	now := time.Now().UTC()

	if u.CreatedAt.After(now) {
		return fmt.Errorf("%w: created at timestamp cannot be in the future", ErrValidation)
	}

	if u.UpdatedAt.After(now) {
		return fmt.Errorf("%w: updated at timestamp cannot be in the future", ErrValidation)
	}

	return nil
}

type UserService interface {
	GetByID(ctx context.Context, userID int64) (*User, error)
	Create(ctx context.Context, favoriteRegions, wantResidency []RegionID) error
	Update(ctx context.Context, userID int64, favoriteRegions, wantResidency []RegionID) error
}

type UserRepository interface {
	GetByID(ctx context.Context, userID int64) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID int64) error
}

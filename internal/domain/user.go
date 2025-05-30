package domain

import (
	"context"
	"fmt"
	"time"
)

var maxRegions = 300

type User struct {
	ID              string    `json:"id"`
	FavoriteRegions []string  `json:"favoriteRegions"`
	WantResidency   []string  `json:"wantResidency"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (u *User) Validate() error {

	if u.ID == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
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
	Create(ctx context.Context, userID string) error
	Update(ctx context.Context, userID string, favoriteRegions, wantResidency []string) error
}

type UserRepository interface {
	GetByID(ctx context.Context, userID string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID string) error
}

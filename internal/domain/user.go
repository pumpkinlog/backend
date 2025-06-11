package domain

import (
	"context"
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
	if u.ID <= 0 {
		return ValidationError("user ID is required")
	}

	if len(u.FavoriteRegions) > maxRegions {
		return ValidationError("favorite regions cannot be greater than %d", maxRegions)
	}

	if len(u.WantResidency) > maxRegions {
		return ValidationError(" want residency cannot be greater than %d", maxRegions)
	}

	if u.CreatedAt.IsZero() {
		return ValidationError("created at timestamp is required")
	}

	if u.UpdatedAt.IsZero() {
		return ValidationError("updated at timestamp is required")
	}

	now := time.Now().UTC()

	if u.CreatedAt.After(now) {
		return ValidationError("created at timestamp cannot be in the future")
	}

	if u.UpdatedAt.After(now) {
		return ValidationError("updated at timestamp cannot be in the future")
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

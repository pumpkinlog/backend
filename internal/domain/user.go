package domain

import (
	"context"
	"time"
)

type User struct {
	ID              string    `json:"id"`
	FavoriteRegions []string  `json:"favoriteRegions"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)

	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error

	Delete(ctx context.Context, id string) error
}

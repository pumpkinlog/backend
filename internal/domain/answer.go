package domain

import (
	"context"
	"fmt"
	"time"
)

type Answer struct {
	UserID      int64     `json:"userId"`
	ConditionID int64     `json:"conditionId"`
	RegionID    string    `json:"regionId"`
	Value       any       `json:"value"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (a *Answer) Validate() error {

	if a.UserID < 0 {
		return fmt.Errorf("%w: user ID is invalid", ErrValidation)
	}

	if a.ConditionID < 0 {
		return fmt.Errorf("%w: condition ID is invalid", ErrValidation)
	}

	if a.Value == nil {
		return fmt.Errorf("%w: value is required", ErrValidation)
	}

	if a.CreatedAt.IsZero() {
		return fmt.Errorf("%w: created at timestamp is invalid", ErrValidation)
	}

	if a.UpdatedAt.IsZero() {
		return fmt.Errorf("%w: updated at timestamp is invalid", ErrValidation)
	}

	now := time.Now().UTC()

	if a.CreatedAt.After(now) {
		return fmt.Errorf("%w: created at timestamp cannot be in the future", ErrValidation)
	}

	if a.UpdatedAt.After(now) {
		return fmt.Errorf("%w: updated at timestamp cannot be in the future", ErrValidation)
	}

	return nil
}

type AnswerService interface {
	CreateOrUpdate(ctx context.Context, userID, conditionID int64, value any) error
}

type AnswerRepository interface {
	GetByID(ctx context.Context, userID, conditionID int64) (*Answer, error)
	ListByRegionID(ctx context.Context, userID int64, regionID string) ([]*Answer, error)
	CreateOrUpdate(ctx context.Context, answer *Answer) error
	Delete(ctx context.Context, userID, conditionID int64) error
}

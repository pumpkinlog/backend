package domain

import (
	"context"
	"time"
)

type Answer struct {
	UserID      int64     `json:"userId"`
	ConditionID Code      `json:"conditionId"`
	RegionID    RegionID  `json:"regionId"`
	Value       any       `json:"value"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (a *Answer) Validate() error {
	if a.UserID <= 0 {
		return ValidationError("user ID is required")
	}

	if err := a.ConditionID.Validate(); err != nil {
		return err
	}

	if err := a.RegionID.Validate(); err != nil {
		return err
	}

	if a.Value == nil {
		return ValidationError("value is required")
	}

	if a.CreatedAt.IsZero() {
		return ValidationError("created at is required")
	}

	if a.UpdatedAt.IsZero() {
		return ValidationError("updated at is required")
	}

	now := time.Now().UTC()

	if a.CreatedAt.After(now) {
		return ValidationError("created at timestamp cannot be in the future")
	}

	if a.UpdatedAt.After(now) {
		return ValidationError("updated at timestamp cannot be in the future")
	}

	return nil
}

type AnswerService interface {
	GetByID(ctx context.Context, userID int64, conditionID Code) (*Answer, error)
	CreateOrUpdate(ctx context.Context, userID int64, conditionID Code, value any) error
	Delete(ctx context.Context, userID int64, conditionID Code) error
}

type AnswerRepository interface {
	GetByID(ctx context.Context, userID int64, conditionID Code) (*Answer, error)
	ListByRegionID(ctx context.Context, userID int64, regionID RegionID) ([]*Answer, error)
	CreateOrUpdate(ctx context.Context, answer *Answer) error
	Delete(ctx context.Context, userID int64, conditionID Code) error
}

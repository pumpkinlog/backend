package domain

import (
	"context"
	"fmt"
	"time"
)

type Answer struct {
	UserID      string    `json:"userId"`
	ConditionID string    `json:"conditionId"`
	Value       any       `json:"value"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (a *Answer) Validate() error {

	if a.UserID == "" {
		return fmt.Errorf("%w: user ID is required", ErrValidation)
	}

	if a.ConditionID == "" {
		return fmt.Errorf("%w: condition ID is required", ErrValidation)
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
	CreateOrUpdate(ctx context.Context, userID, conditionID string, value any) error
}

type AnswerFilter struct {
	ConditionIDs []string
}

type AnswerRepository interface {
	GetByID(ctx context.Context, userID, conditionID string) (*Answer, error)
	List(ctx context.Context, userID string, filter *AnswerFilter) ([]*Answer, error)
	CreateOrUpdate(ctx context.Context, answer *Answer) error
	Delete(ctx context.Context, userID, conditionID string) error
}

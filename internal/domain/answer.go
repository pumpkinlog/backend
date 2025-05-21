package domain

import "context"

type Answer struct {
	UserID      string `json:"user_Id"`
	ConditionID string `json:"conditionId"`
	Value       any    `json:"value"`
}

type AnswerFilter struct {
	ConditionIDs []string
}

type AnswerRepository interface {
	GetByID(ctx context.Context, userID, conditionID string) (*Answer, error)
	List(ctx context.Context, userID string, filter *AnswerFilter) ([]*Answer, error)

	CreateOrUpdate(ctx context.Context, answer *Answer) error
}

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresAnswerRepository struct {
	conn Connection
}

func NewPostgresAnswerRepository(conn Connection) domain.AnswerRepository {
	return &postgresAnswerRepository{conn}
}

func (r *postgresAnswerRepository) fetch(ctx context.Context, query string, args ...any) ([]*domain.Answer, error) {

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	answers := make([]*domain.Answer, 0)

	for rows.Next() {
		var answer domain.Answer
		if err := rows.Scan(
			&answer.UserID,
			&answer.ConditionID,
			&answer.Value,
			&answer.CreatedAt,
			&answer.UpdatedAt,
		); err != nil {
			return nil, err
		}
		answers = append(answers, &answer)
	}

	return answers, nil
}

func (r *postgresAnswerRepository) GetByID(ctx context.Context, userID, conditionID string) (*domain.Answer, error) {

	query := `
			SELECT 
				user_id,
				condition_id,
				value,
				created_at,
				updated_at
			FROM answers
			WHERE user_id = $1 AND condition_id = $2`

	answers, err := r.fetch(ctx, query, userID, conditionID)
	if err != nil {
		return nil, err
	}

	if len(answers) == 0 {
		return nil, domain.ErrNotFound
	}

	return answers[0], nil
}

func (r *postgresAnswerRepository) List(ctx context.Context, userID string, filter *domain.AnswerFilter) ([]*domain.Answer, error) {
	var query strings.Builder
	args := []any{userID}
	argIndex := 2

	query.WriteString(`
		SELECT 
			user_id,
			condition_id,
			value,
			created_at,
			updated_at
		FROM answers
		WHERE user_id = $1
	`)

	query.WriteString(" AND 1=1")

	if len(filter.ConditionIDs) > 0 {
		query.WriteString(" AND condition_id IN (")
		for i, id := range filter.ConditionIDs {
			if i > 0 {
				query.WriteString(", ")
			}
			query.WriteString(fmt.Sprintf("$%d", argIndex))
			args = append(args, id)
			argIndex++
		}
		query.WriteString(")")
	}

	return r.fetch(ctx, query.String(), args...)
}

func (r *postgresAnswerRepository) CreateOrUpdate(ctx context.Context, answer *domain.Answer) error {

	if answer == nil {
		answer = &domain.Answer{}
	}

	query := `
			INSERT INTO answers (user_id, condition_id, value, created_at, updated_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, condition_id)
			DO UPDATE SET value = $3, updated_at = $5`

	_, err := r.conn.Exec(
		ctx,
		query,
		answer.UserID,
		answer.ConditionID,
		answer.Value,
		answer.CreatedAt,
		answer.UpdatedAt,
	)
	return err
}

func (r *postgresAnswerRepository) Delete(ctx context.Context, userID, conditionID string) error {

	query := `
		DELETE FROM answers
		WHERE user_id = $1 AND condition_id = $2`

	_, err := r.conn.Exec(ctx, query, userID, conditionID)
	return err
}

package repository

import (
	"context"
	"errors"

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
			&answer.RegionID,
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

func (r *postgresAnswerRepository) GetByID(ctx context.Context, userID int64, conditionID domain.Code) (*domain.Answer, error) {

	query := `
			SELECT 
				user_id,
				condition_id,
				region_id,
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

func (r *postgresAnswerRepository) ListByRegionID(ctx context.Context, userID int64, regionID domain.RegionID) ([]*domain.Answer, error) {

	query := `
		SELECT 
			user_id,
			condition_id,
			region_id,
			value,
			created_at,
			updated_at
		FROM answers
		WHERE user_id = $1 AND region_id = $2`

	return r.fetch(ctx, query, userID, regionID)
}

func (r *postgresAnswerRepository) CreateOrUpdate(ctx context.Context, answer *domain.Answer) error {

	if answer == nil {
		return errors.New("answer cannot be nil")
	}

	query := `
			INSERT INTO answers (user_id, condition_id, region_id, value, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (user_id, condition_id)
			DO UPDATE SET value = $4, updated_at = $5`

	_, err := r.conn.Exec(
		ctx,
		query,
		answer.UserID,
		answer.ConditionID,
		answer.RegionID,
		answer.Value,
		answer.CreatedAt,
		answer.UpdatedAt,
	)
	return err
}

func (r *postgresAnswerRepository) Delete(ctx context.Context, userID int64, conditionID domain.Code) error {

	query := `
		DELETE FROM answers
		WHERE user_id = $1 AND condition_id = $2`

	_, err := r.conn.Exec(ctx, query, userID, conditionID)
	return err
}

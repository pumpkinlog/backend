package repository

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresEvaluationRepository struct {
	conn Connection
}

func NewPostgresEvaluationRepository(conn Connection) domain.EvaluationRepository {
	return &postgresEvaluationRepository{conn}
}

func (r *postgresEvaluationRepository) fetch(ctx context.Context, query string, args ...any) ([]*domain.RegionEvaluation, error) {

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	evaluations := make([]*domain.RegionEvaluation, 0)

	for rows.Next() {
		var evaluation domain.RegionEvaluation
		if err := rows.Scan(
			&evaluation.UserID,
			&evaluation.RegionID,
			&evaluation.Passed,
			&evaluation.RuleEvaluations,
			&evaluation.EvaluatedAt,
		); err != nil {
			return nil, err
		}
		evaluations = append(evaluations, &evaluation)
	}

	return evaluations, nil
}

func (r *postgresEvaluationRepository) GetByID(ctx context.Context, userID int64, regionID string) (*domain.RegionEvaluation, error) {

	query := `
		SELECT 
			user_id,
			region_id,
			passed,
			evaluations,
			evaluated_at
		FROM evaluations
		WHERE user_id = $1 AND region_id = $2`

	evaluations, err := r.fetch(ctx, query, userID, regionID)
	if err != nil {
		return nil, err
	}

	if len(evaluations) == 0 {
		return nil, domain.ErrNotFound
	}

	return evaluations[0], nil
}

func (r *postgresEvaluationRepository) List(ctx context.Context, userID int64) ([]*domain.RegionEvaluation, error) {

	query := `
		SELECT 
			user_id,
			region_id,
			passed,
			evaluations,
			evaluated_at
		FROM evaluations
		WHERE user_id = $1`

	evaluations, err := r.fetch(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	return evaluations, nil
}

func (r *postgresEvaluationRepository) CreateOrUpdate(ctx context.Context, evaluation *domain.RegionEvaluation) error {

	query := `
		INSERT INTO evaluations (user_id, region_id, passed, evaluations, evaluated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, region_id) DO UPDATE
		SET passed = $3, evaluations = $4, evaluated_at = $5`

	_, err := r.conn.Exec(
		ctx,
		query,
		evaluation.UserID,
		evaluation.RegionID,
		evaluation.Passed,
		evaluation.RuleEvaluations,
		evaluation.EvaluatedAt,
	)
	return err
}

func (r *postgresEvaluationRepository) DeleteByRegionID(ctx context.Context, regionID string) error {

	query := `
		DELETE FROM evaluations
		WHERE region_id = $1`

	_, err := r.conn.Exec(ctx, query, regionID)
	return err
}

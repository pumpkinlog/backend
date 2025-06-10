package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresConditionRepository struct {
	conn Connection
}

func NewPostgresConditionRepository(conn Connection) domain.ConditionRepository {
	return &postgresConditionRepository{conn}
}

func (r *postgresConditionRepository) fetch(ctx context.Context, query string, args ...any) ([]*domain.Condition, error) {

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conditions := make([]*domain.Condition, 0)

	for rows.Next() {
		var condition domain.Condition
		if err := rows.Scan(
			&condition.ID,
			&condition.RegionID,
			&condition.Prompt,
			&condition.Type,
		); err != nil {
			return nil, err
		}
		conditions = append(conditions, &condition)
	}

	return conditions, nil
}

func (r *postgresConditionRepository) GetByID(ctx context.Context, conditionID domain.Code) (*domain.Condition, error) {

	query := `
		SELECT
			id,
			region_id,
			prompt,
			type
		FROM conditions
		WHERE id = $1`

	conditions, err := r.fetch(ctx, query, conditionID)
	if err != nil {
		return nil, err
	}

	if len(conditions) == 0 {
		return nil, domain.ErrNotFound
	}

	return conditions[0], nil
}

func (r *postgresConditionRepository) List(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error) {
	var (
		query strings.Builder
		args  []any
	)

	query.WriteString(`
		SELECT
			id,
			region_id,
			prompt,
			type
		FROM conditions
		GROUP BY id, region_id
		ORDER BY id`)
	argIndex := 1

	if len(filter.RegionIDs) > 0 {
		query.WriteString(" WHERE region_id IN (")
		for i, id := range filter.RegionIDs {
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

func (r *postgresConditionRepository) ListByRegionID(ctx context.Context, regionID domain.RegionID) ([]*domain.Condition, error) {

	query := `
		SELECT
			id,
			region_id,
			prompt,
			type
		FROM conditions
		WHERE region_id = $1
		GROUP BY id, region_id
		ORDER BY id`

	return r.fetch(ctx, query, regionID)
}

func (r *postgresConditionRepository) CreateOrUpdate(ctx context.Context, condition *domain.Condition) error {

	query := `
			INSERT INTO conditions (
				id,
				region_id,
				prompt,
				type
			) VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET
				prompt = $3,
				type = $4`

	_, err := r.conn.Exec(
		ctx,
		query,
		condition.ID,
		condition.RegionID,
		condition.Prompt,
		condition.Type,
	)
	return err
}

func (r *postgresConditionRepository) Delete(ctx context.Context, conditionID domain.Code) error {

	query := `
			DELETE FROM conditions
			WHERE id = $1`

	_, err := r.conn.Exec(ctx, query, conditionID)
	return err
}

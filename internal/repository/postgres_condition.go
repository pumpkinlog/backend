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
			&condition.Prompt,
			&condition.Type,
		); err != nil {
			return nil, err
		}
		conditions = append(conditions, &condition)
	}

	return conditions, nil
}

func (r *postgresConditionRepository) GetByID(ctx context.Context, conditionID string) (*domain.Condition, error) {

	query := `
			SELECT id, prompt, type
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

	if filter == nil {
		filter = new(domain.ConditionFilter)
	}

	var (
		query strings.Builder
		args  []any
	)

	query.WriteString(`
		SELECT id, prompt, type
		FROM conditions
		WHERE 1=1
	`)
	argIndex := 1

	if len(filter.ConditionIDs) > 0 {
		query.WriteString(" AND id IN (")
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

	if filter.Limit != nil && filter.Page != nil && *filter.Limit > 0 && *filter.Page > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT $%d", argIndex))
		args = append(args, *filter.Limit)
		argIndex++

		offset := (*filter.Page - 1) * (*filter.Limit)
		query.WriteString(fmt.Sprintf(" OFFSET $%d", argIndex))
		args = append(args, offset)
	}

	return r.fetch(ctx, query.String(), args...)
}

func (r *postgresConditionRepository) CreateOrUpdate(ctx context.Context, condition *domain.Condition) error {

	query := `
			INSERT INTO conditions (
				id,
				prompt,
				type
			) VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE SET
				prompt = $2,
				type = $3`

	_, err := r.conn.Exec(
		ctx,
		query,
		condition.ID,
		condition.Prompt,
		condition.Type,
	)
	return err
}

func (r *postgresConditionRepository) Delete(ctx context.Context, conditionID string) error {

	query := `
			DELETE FROM conditions
			WHERE id = $1`

	_, err := r.conn.Exec(ctx, query, conditionID)
	return err
}

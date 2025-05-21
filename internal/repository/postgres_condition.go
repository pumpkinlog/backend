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
			&condition.RuleID,
			&condition.Prompt,
			&condition.Type,
			&condition.Comparator,
			&condition.Expected,
		); err != nil {
			return nil, err
		}
		conditions = append(conditions, &condition)
	}

	return conditions, nil
}

func (r *postgresConditionRepository) GetByID(ctx context.Context, id string) (*domain.Condition, error) {

	query := `
			SELECT id, rule_id, prompt, type, comparator, expected
			FROM conditions
			WHERE id = $1`

	conditions, err := r.fetch(ctx, query, id)
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
		SELECT id, rule_id, prompt, type, comparator, expected
		FROM conditions
		WHERE 1=1
	`)
	argIndex := 1

	if len(filter.RuleIDs) > 0 {
		query.WriteString(" AND rule_id IN (")
		for i, id := range filter.RuleIDs {
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
				rule_id, 
				prompt,
				type,
				comparator,
				expected
			) VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE SET
				rule_id = $2,
				prompt = $3,
				type = $4,
				comparator = $5,
				expected = $6`

	_, err := r.conn.Exec(
		ctx,
		query,
		condition.ID,
		condition.RuleID,
		condition.Prompt,
		condition.Type,
		condition.Comparator,
		condition.Expected,
	)
	return err
}

func (r *postgresConditionRepository) Delete(ctx context.Context, id string) error {

	query := `
			DELETE FROM conditions
			WHERE id = $1`

	_, err := r.conn.Exec(ctx, query, id)
	return err
}

package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresRuleConditionRepository struct {
	conn Connection
}

func NewPostgresRuleConditionRepository(conn Connection) domain.RuleConditionRepository {
	return &postgresRuleConditionRepository{conn}
}

func (r *postgresRuleConditionRepository) fetch(ctx context.Context, query string, args ...any) ([]*domain.RuleCondition, error) {
	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conditions := make([]*domain.RuleCondition, 0)

	for rows.Next() {
		var condition domain.RuleCondition
		if err := rows.Scan(
			&condition.RuleID,
			&condition.ConditionID,
			&condition.Weight,
			&condition.Comparator,
			&condition.Expected,
		); err != nil {
			return nil, err
		}
		conditions = append(conditions, &condition)
	}

	return conditions, nil
}

func (r *postgresRuleConditionRepository) GetByID(ctx context.Context, ruleID, conditionID string) (*domain.RuleCondition, error) {

	query := `
		SELECT 
			rule_id,
			condition_id,
			weight,
			comparator,
			expected
		FROM rule_conditions
		WHERE rule_id = $1 AND condition_id = $2`

	conditions, err := r.fetch(ctx, query, ruleID, conditionID)
	if err != nil {
		return nil, err
	}

	if len(conditions) == 0 {
		return nil, domain.ErrNotFound
	}

	return conditions[0], nil
}

func (r *postgresRuleConditionRepository) List(ctx context.Context, filter *domain.RuleConditionFilter) ([]*domain.RuleCondition, error) {
	if filter == nil {
		filter = new(domain.RuleConditionFilter)
	}

	var (
		query    strings.Builder
		args     []any
		argIndex = 1
	)

	query.WriteString(`
		SELECT 
			rule_id,
			condition_id,
			weight,
			comparator,
			expected
		FROM rule_conditions
		WHERE 1=1
	`)

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

	query.WriteString(" ORDER BY rule_id, condition_id")

	return r.fetch(ctx, query.String(), args...)
}

func (r *postgresRuleConditionRepository) CreateOrUpdate(ctx context.Context, rc *domain.RuleCondition) error {

	if rc == nil {
		rc = new(domain.RuleCondition)
	}

	query := `
		INSERT INTO rule_conditions (
			rule_id,
			condition_id,
			weight,
			comparator,
			expected
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (rule_id, condition_id) DO UPDATE SET
			weight = $3,
			comparator = $4,
			expected = $5`

	_, err := r.conn.Exec(
		ctx,
		query,
		rc.RuleID,
		rc.ConditionID,
		rc.Weight,
		rc.Comparator,
		rc.Expected,
	)
	return err
}

func (r *postgresRuleConditionRepository) Delete(ctx context.Context, ruleID, conditionID string) error {
	query := `
		DELETE FROM rule_conditions
		WHERE rule_id = $1 AND condition_id = $2`

	_, err := r.conn.Exec(ctx, query, ruleID, conditionID)
	return err
}

package repository

import (
	"context"

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
			&condition.RegionID,
			&condition.Comparator,
			&condition.Expected,
		); err != nil {
			return nil, err
		}
		conditions = append(conditions, &condition)
	}

	return conditions, nil
}

func (r *postgresRuleConditionRepository) GetByID(ctx context.Context, ruleID, conditionID int64) (*domain.RuleCondition, error) {

	query := `
		SELECT 
			rule_id,
			condition_id,
			region_id,
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

func (r *postgresRuleConditionRepository) ListByRegionID(ctx context.Context, regionID string) ([]*domain.RuleCondition, error) {

	query := `
		SELECT 
			rule_id,
			condition_id,
			region_id,
			comparator,
			expected
		FROM rule_conditions
		WHERE region_id = $1`

	return r.fetch(ctx, query, regionID)
}

func (r *postgresRuleConditionRepository) CreateOrUpdate(ctx context.Context, rc *domain.RuleCondition) error {

	if rc == nil {
		rc = new(domain.RuleCondition)
	}

	query := `
		INSERT INTO rule_conditions (
			rule_id,
			condition_id,
			region_id,
			comparator,
			expected
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (rule_id, condition_id) DO UPDATE SET
			comparator = $4,
			expected = $5`

	_, err := r.conn.Exec(
		ctx,
		query,
		rc.RuleID,
		rc.ConditionID,
		rc.RegionID,
		rc.Comparator,
		rc.Expected,
	)
	return err
}

func (r *postgresRuleConditionRepository) Delete(ctx context.Context, ruleID, conditionID int64) error {
	query := `
		DELETE FROM rule_conditions
		WHERE rule_id = $1 AND condition_id = $2`

	_, err := r.conn.Exec(ctx, query, ruleID, conditionID)
	return err
}

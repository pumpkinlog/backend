package repository

import (
	"context"
	"fmt"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresRuleConditionRepository struct {
	conn Connection
}

func NewPostgresRuleConditionRepository(conn Connection) domain.RuleConditionRepository {
	return &postgresRuleConditionRepository{conn}
}

func (r *postgresRuleConditionRepository) Link(ctx context.Context, ruleID, conditionID string) error {

	query := `
		INSERT INTO rule_conditions (rule_id, condition_id)
		VALUES ($1, $2)
		ON CONFLICT (rule_id, condition_id) DO NOTHING
	`

	_, err := r.conn.Exec(ctx, query, ruleID, conditionID)
	if err != nil {
		return fmt.Errorf("failed to link rule to condition: %w", err)
	}

	return nil
}

func (r *postgresRuleConditionRepository) Unlink(ctx context.Context, ruleID, conditionID string) error {

	query := `
		DELETE FROM rule_conditions
		WHERE rule_id = $1 AND condition_id = $2
	`

	_, err := r.conn.Exec(ctx, query, ruleID, conditionID)
	if err != nil {
		return fmt.Errorf("failed to unlink rule from condition: %w", err)
	}

	return nil
}

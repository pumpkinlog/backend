package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresRuleRepository struct {
	conn Connection
}

func NewPostgresRuleRepository(conn Connection) domain.RuleRepository {
	return &postgresRuleRepository{conn}
}

func (p *postgresRuleRepository) fetch(ctx context.Context, query string, args ...any) ([]*domain.Rule, error) {

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := make([]*domain.Rule, 0)

	for rows.Next() {
		var rule domain.Rule
		if err := rows.Scan(
			&rule.ID,
			&rule.RegionID,
			&rule.Name,
			&rule.Description,
			&rule.Node,
		); err != nil {
			return nil, err
		}

		rules = append(rules, &rule)
	}

	return rules, nil
}

func (r *postgresRuleRepository) GetByID(ctx context.Context, ruleID domain.Code) (*domain.Rule, error) {

	query := `SELECT
				id,
				region_id,
				name,
				description,
				node
			FROM rules WHERE id = $1`

	rules, err := r.fetch(ctx, query, ruleID)
	if err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return nil, domain.ErrNotFound
	}

	return rules[0], nil
}

func (r *postgresRuleRepository) List(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error) {

	if filter == nil {
		filter = new(domain.RuleFilter)
	}

	var (
		query strings.Builder
		args  []any
	)

	query.WriteString(`
		SELECT 
			id,
			region_id,
			name,
			description,
			node
		FROM rules
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

func (r *postgresRuleRepository) ListByRegionID(ctx context.Context, regionID domain.RegionID) ([]*domain.Rule, error) {

	query := `
		SELECT 
				id,
				region_id,
				name,
				description,
				node
		FROM rules
		WHERE region_id = $1
		GROUP BY id, region_id
		ORDER BY id`

	return r.fetch(ctx, query, regionID)
}

func (r *postgresRuleRepository) CreateOrUpdate(ctx context.Context, rule *domain.Rule) error {

	query := `
			INSERT INTO rules (
				id,
				region_id,
				name,
				description,
				node
			) VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO UPDATE SET
				name = $3,
				description = $4,
				node = $5`

	_, err := r.conn.Exec(
		ctx,
		query,
		rule.ID,
		rule.RegionID,
		rule.Name,
		rule.Description,
		rule.Node,
	)
	return err
}

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
			&rule.RuleType,
			&rule.PeriodType,
			&rule.Threshold,
			&rule.YearStartMonth,
			&rule.YearStartDay,
			&rule.OffsetYears,
			&rule.Years,
			&rule.RollingDays,
			&rule.RollingMonths,
			&rule.RollingYears,
		); err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}

func (r *postgresRuleRepository) GetByID(ctx context.Context, id string) (*domain.Rule, error) {

	query := `SELECT
				id,
				region_id,
				name,
				description,
				rule_type,
				period_type,
				threshold,
				year_start_month,
				year_start_day,
				offset_years,
				years,
				rolling_days,
				rolling_months,
				rolling_years
			FROM rules WHERE id = $1`

	rules, err := r.fetch(ctx, query, id)
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
				rule_type,
				period_type,
				threshold,
				year_start_month,
				year_start_day,
				offset_years,
				years,
				rolling_days,
				rolling_months,
				rolling_years
		FROM rules
		WHERE 1=1
	`)
	argIndex := 1

	if len(filter.RegionIDs) > 0 {
		query.WriteString(" AND region_id IN (")
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

func (r *postgresRuleRepository) CreateOrUpdate(ctx context.Context, rule *domain.Rule) error {

	query := `
			INSERT INTO rules (
				id,
				region_id,
				name,
				description,
				rule_type,
				period_type,
				threshold,
				year_start_month,
				year_start_day,
				offset_years,
				years,
				rolling_days,
				rolling_months,
				rolling_years
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			ON CONFLICT (id) DO UPDATE SET
				region_id = $2,
				name = $3,
				description = $4,
				rule_type = $5,
				period_type = $6,
				threshold = $7,
				year_start_month = $8,
				year_start_day = $9,
				offset_years = $10,
				years = $11,
				rolling_days = $12,
				rolling_months = $13,
				rolling_years = $14`

	_, err := r.conn.Exec(
		ctx,
		query,
		rule.ID,
		rule.RegionID,
		rule.Name,
		rule.Description,
		rule.RuleType,
		rule.PeriodType,
		rule.Threshold,
		rule.YearStartMonth,
		rule.YearStartDay,
		rule.OffsetYears,
		rule.Years,
		rule.RollingDays,
		rule.RollingMonths,
		rule.RollingYears,
	)
	return err
}

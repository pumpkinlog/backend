package repository

import (
	"context"
	"strings"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresRegionConditionRepository struct {
	conn Connection
}

func NewPostgresRegionConditionRepository(conn Connection) domain.RegionConditionRepository {
	return &postgresRegionConditionRepository{conn}
}

func (r *postgresRegionConditionRepository) fetch(ctx context.Context, query string, args ...any) ([]*domain.RegionCondition, error) {
	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*domain.RegionCondition, 0)
	for rows.Next() {
		var rc domain.RegionCondition
		if err := rows.Scan(&rc.RegionID, &rc.ConditionID); err != nil {
			return nil, err
		}
		results = append(results, &rc)
	}
	return results, nil
}

func (r *postgresRegionConditionRepository) GetByID(ctx context.Context, regionID, conditionID string) (*domain.RegionCondition, error) {

	query := `
		SELECT region_id, condition_id
		FROM region_conditions
		WHERE region_id = $1 AND condition_id = $2`

	results, err := r.fetch(ctx, query, regionID, conditionID)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, domain.ErrNotFound
	}

	return results[0], nil
}

func (r *postgresRegionConditionRepository) List(ctx context.Context, filter *domain.RegionConditionFilter) ([]*domain.RegionCondition, error) {

	var (
		query strings.Builder
		args  []any
	)

	query.WriteString(`SELECT region_id, condition_id FROM region_conditions`)

	if filter != nil && filter.Page > 0 && filter.Limit > 0 {
		query.WriteString(" LIMIT $1 OFFSET $2")
		args = append(args, filter.Limit, (filter.Page-1)*filter.Limit)
	}

	return r.fetch(ctx, query.String(), args...)
}

func (r *postgresRegionConditionRepository) CreateOrUpdate(ctx context.Context, rc *domain.RegionCondition) error {

	if rc == nil {
		rc = new(domain.RegionCondition)
	}

	query := `
		INSERT INTO region_conditions (region_id, condition_id)
		VALUES ($1, $2)
		ON CONFLICT (region_id, condition_id) DO NOTHING`

	_, err := r.conn.Exec(ctx, query, rc.RegionID, rc.ConditionID)
	return err
}

func (r *postgresRegionConditionRepository) Delete(ctx context.Context, regionID, conditionID string) error {

	query := `
		DELETE FROM region_conditions
		WHERE region_id = $1 AND condition_id = $2`

	_, err := r.conn.Exec(ctx, query, regionID, conditionID)
	return err
}

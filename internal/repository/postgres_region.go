package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresRegionRepository struct {
	conn Connection
}

func NewPostgresRegionRepository(conn Connection) domain.RegionRepository {
	return &postgresRegionRepository{conn}
}

func (p *postgresRegionRepository) fetch(ctx context.Context, query string, args ...any) ([]*domain.Region, error) {

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	regions := make([]*domain.Region, 0)

	for rows.Next() {
		var region domain.Region
		if err := rows.Scan(
			&region.ID,
			&region.ParentRegionID,
			&region.Type,
			&region.Name,
			&region.Continent,
			&region.YearStartMonth,
			&region.YearStartDay,
			&region.LatLng,
		); err != nil {
			return nil, err
		}
		regions = append(regions, &region)
	}

	return regions, nil
}

func (r *postgresRegionRepository) GetByID(ctx context.Context, id string) (*domain.Region, error) {

	query := `
		SELECT 
			id, 
			parent_region_id, 
			region_type, 
			name, 
			continent, 
			year_start_month,
			year_start_day,
			lat_lng
		FROM regions
		WHERE id = $1`

	regions, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(regions) == 0 {
		return nil, domain.ErrNotFound
	}

	return regions[0], nil
}

func (r *postgresRegionRepository) List(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error) {

	if filter == nil {
		filter = new(domain.RegionFilter)
	}

	var (
		query    strings.Builder
		args     []any
		argIndex = 1
	)

	query.WriteString(`
		SELECT 
			id, 
			parent_region_id, 
			region_type, 
			name, 
			continent, 
			year_start_month,
			year_start_day,
			lat_lng
		FROM regions
		WHERE 1=1
	`)

	if len(filter.RegionIDs) > 0 {
		query.WriteString(" AND id IN (")
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

	query.WriteString(" ORDER BY id")

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

func (r *postgresRegionRepository) CreateOrUpdate(ctx context.Context, region *domain.Region) error {

	query := `
		INSERT INTO regions (
			id, 
			parent_region_id, 
			region_type, 
			name, 
			continent, 
			year_start_month,
			year_start_day,
			lat_lng
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			parent_region_id = $2,
			region_type = $3,
			name = $4,
			continent = $5,
			year_start_month = $6,
			year_start_day = $7,
			lat_lng = $8`

	_, err := r.conn.Exec(ctx, query,
		region.ID,
		region.ParentRegionID,
		region.Type,
		region.Name,
		region.Continent,
		region.YearStartMonth,
		region.YearStartDay,
		region.LatLng,
	)
	return err
}

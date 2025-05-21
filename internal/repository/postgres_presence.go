package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresPresenceRepository struct {
	conn Connection
}

func NewPostgresPresenceRepository(conn Connection) domain.PresenceRepository {
	return &postgresPresenceRepository{conn}
}

func (p *postgresPresenceRepository) fetch(ctx context.Context, query string, args ...any) ([]*domain.Presence, error) {

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	presences := make([]*domain.Presence, 0)

	for rows.Next() {
		var presence domain.Presence
		if err := rows.Scan(
			&presence.UserID,
			&presence.RegionID,
			&presence.Date,
			&presence.DeviceID,
		); err != nil {
			return nil, err
		}
		presences = append(presences, &presence)
	}

	return presences, nil
}

func (r *postgresPresenceRepository) GetByID(ctx context.Context, userID, regionID string, date time.Time) (*domain.Presence, error) {

	query := `
			SELECT user_id, region_id, date, device_id
			FROM presences
			WHERE user_id = $1 AND region_id = $2 AND date = $3`

	locations, err := r.fetch(ctx, query, userID, regionID, date)
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, domain.ErrNotFound
	}

	return locations[0], nil
}

func (r *postgresPresenceRepository) List(ctx context.Context, userID string, filter *domain.PresenceFilter) ([]*domain.Presence, error) {

	if filter == nil {
		filter = new(domain.PresenceFilter)
	}

	var (
		query    strings.Builder
		args     []any
		argIndex = 1
	)

	query.WriteString(`
		SELECT 
			user_id,
			region_id,
			date,
			device_id
		FROM presences
		WHERE user_id = $1
	`)
	args = append(args, userID)
	argIndex++

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

	if filter.Start != nil && filter.End != nil {
		query.WriteString(fmt.Sprintf(" AND date BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, *filter.Start, *filter.End)
		argIndex += 2
	} else if filter.Start != nil {
		query.WriteString(fmt.Sprintf(" AND date >= $%d", argIndex))
		args = append(args, *filter.Start)
		argIndex++
	} else if filter.End != nil {
		query.WriteString(fmt.Sprintf(" AND date <= $%d", argIndex))
		args = append(args, *filter.End)
		argIndex++
	}

	query.WriteString(" ORDER BY date")

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

func (r *postgresPresenceRepository) ListByRegionBounds(ctx context.Context, userID string, bounds map[string]domain.TimeWindow) ([]*domain.Presence, error) {

	if len(bounds) == 0 {
		return nil, errors.New("no region bounds provided")
	}

	var (
		query         strings.Builder
		args          []any
		argIndex      = 1
		regionClauses []string
	)

	query.WriteString(`
		SELECT user_id, region_id, date, device_id
		FROM presences
		WHERE user_id = $1 AND (
	`)
	args = append(args, userID)
	argIndex++

	for regionID, window := range bounds {
		regionClauses = append(regionClauses, fmt.Sprintf(
			"(region_id = $%d AND date BETWEEN $%d AND $%d)",
			argIndex, argIndex+1, argIndex+2,
		))
		args = append(args, regionID, window.Start, window.End)
		argIndex += 3
	}

	query.WriteString(strings.Join(regionClauses, " OR "))
	query.WriteString(") ORDER BY region_id, date")

	return r.fetch(ctx, query.String(), args...)
}

func (r *postgresPresenceRepository) Create(ctx context.Context, location *domain.Presence) error {

	if location == nil {
		location = &domain.Presence{}
	}

	query := `
			INSERT INTO presences (user_id, region_id, date, device_id)
			VALUES ($1, $2, $3, $4)`

	_, err := r.conn.Exec(
		ctx,
		query,
		location.UserID,
		location.RegionID,
		location.Date,
		location.DeviceID,
	)
	return err
}

func (r *postgresPresenceRepository) CreateRange(ctx context.Context, userID, regionID string, deviceID *string, start, end time.Time) error {

	query := `
			INSERT INTO presences (user_id, region_id, date, device_id)
			SELECT $1, $2, d::date, $3
			FROM generate_series($4::date, $5::date, '1 day') AS d`

	_, err := r.conn.Exec(ctx, query, userID, regionID, deviceID, start, end)
	return err
}

func (r *postgresPresenceRepository) Delete(ctx context.Context, userID, regionID string, date time.Time) error {

	query := `
			DELETE FROM presences
			WHERE user_id = $1 AND region_id = $2 AND date = $3`

	_, err := r.conn.Exec(ctx, query, userID, regionID, date)
	return err
}

func (r *postgresPresenceRepository) DeleteRange(ctx context.Context, userID, regionID string, start, end time.Time) error {

	query := `
			DELETE FROM presences
			WHERE user_id = $1
			AND region_id = $2 
			AND date BETWEEN $3 AND $4`

	_, err := r.conn.Exec(ctx, query, userID, regionID, start, end)
	return err
}

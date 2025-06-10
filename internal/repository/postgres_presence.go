package repository

import (
	"context"
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
			&presence.CreatedAt,
			&presence.UpdatedAt,
		); err != nil {
			return nil, err
		}
		presences = append(presences, &presence)
	}

	return presences, nil
}

func (r *postgresPresenceRepository) GetByID(ctx context.Context, userID int64, regionID domain.RegionID, date time.Time) (*domain.Presence, error) {

	query := `
			SELECT user_id, region_id, date, device_id, created_at, updated_at
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

func (r *postgresPresenceRepository) List(ctx context.Context, userID int64, filter *domain.PresenceFilter) ([]*domain.Presence, error) {
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
			device_id,
			created_at,
			updated_at
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
	} else if filter.Start != nil {
		query.WriteString(fmt.Sprintf(" AND date >= $%d", argIndex))
		args = append(args, *filter.Start)
	} else if filter.End != nil {
		query.WriteString(fmt.Sprintf(" AND date <= $%d", argIndex))
		args = append(args, *filter.End)
	}

	query.WriteString(" ORDER BY date")

	return r.fetch(ctx, query.String(), args...)
}

func (r *postgresPresenceRepository) ListByRegionPeriod(ctx context.Context, userID int64, regionID domain.RegionID, start, end time.Time) ([]*domain.Presence, error) {

	query := `
		SELECT user_id, region_id, date, device_id, created_at, updated_at
		FROM presences
		WHERE user_id = $1 AND region_id = $2 AND date BETWEEN $3 AND $4
		ORDER BY date`

	return r.fetch(ctx, query, userID, regionID, start, end)
}

func (r *postgresPresenceRepository) Create(ctx context.Context, presence *domain.Presence) error {

	if presence == nil {
		presence = &domain.Presence{}
	}

	query := `
			INSERT INTO presences (user_id, region_id, date, device_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.conn.Exec(
		ctx,
		query,
		presence.UserID,
		presence.RegionID,
		presence.Date,
		presence.DeviceID,
		presence.CreatedAt,
		presence.UpdatedAt,
	)
	return err
}

func (r *postgresPresenceRepository) CreateRange(ctx context.Context, userID int64, regionID domain.RegionID, deviceID *int64, start, end time.Time) error {

	query := `
			INSERT INTO presences (user_id, region_id, date, device_id, created_at, updated_at)
			SELECT $1, $2, d::date, $3, $6, $7
			FROM generate_series($4::date, $5::date, '1 day') AS d
            ON CONFLICT (user_id, region_id, date) DO NOTHING`

	now := time.Now().UTC()

	_, err := r.conn.Exec(ctx, query, userID, regionID, deviceID, start, end, now, now)
	return err
}

func (r *postgresPresenceRepository) Delete(ctx context.Context, userID int64, regionID domain.RegionID, date time.Time) error {

	query := `
			DELETE FROM presences
			WHERE user_id = $1 AND region_id = $2 AND date = $3`

	_, err := r.conn.Exec(ctx, query, userID, regionID, date)
	return err
}

func (r *postgresPresenceRepository) DeleteRange(ctx context.Context, userID int64, regionID domain.RegionID, start, end time.Time) error {

	query := `
			DELETE FROM presences
			WHERE user_id = $1
			AND region_id = $2 
			AND date BETWEEN $3 AND $4`

	_, err := r.conn.Exec(ctx, query, userID, regionID, start, end)
	return err
}

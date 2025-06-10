package repository

import (
	"context"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresDeviceRepository struct {
	conn Connection
}

func NewPostgresDeviceRepository(conn Connection) domain.DeviceRepository {
	return &postgresDeviceRepository{conn}
}

func (r *postgresDeviceRepository) fetch(ctx context.Context, query string, args ...any) ([]*domain.Device, error) {

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := make([]*domain.Device, 0)

	for rows.Next() {
		var device domain.Device
		if err := rows.Scan(
			&device.ID,
			&device.UserID,
			&device.Name,
			&device.Platform,
			&device.Model,
			&device.Token,
			&device.Active,
			&device.CreatedAt,
			&device.UpdatedAt,
		); err != nil {
			return nil, err
		}
		devices = append(devices, &device)
	}

	return devices, nil
}

func (r *postgresDeviceRepository) GetByID(ctx context.Context, userID, deviceID int64) (*domain.Device, error) {

	query := `
			SELECT id, user_id, name, platform, model, token, active, created_at, updated_at
			FROM devices
			WHERE id = $1 AND user_id = $2`

	devices, err := r.fetch(ctx, query, deviceID, userID)
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, domain.ErrNotFound
	}

	return devices[0], nil
}

func (r *postgresDeviceRepository) List(ctx context.Context, userID int64) ([]*domain.Device, error) {

	query := `
			SELECT id, user_id, name, platform, model, token, active, created_at, updated_at
			FROM devices
			WHERE user_id = $1`

	return r.fetch(ctx, query, userID)
}

func (r *postgresDeviceRepository) Create(ctx context.Context, device *domain.Device) error {

	if device == nil {
		device = &domain.Device{}
	}

	query := `
			INSERT INTO devices (user_id, name, platform, model, token, active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`

	return r.conn.QueryRow(
		ctx,
		query,
		device.UserID,
		device.Name,
		device.Platform,
		device.Model,
		device.Token,
		device.Active,
		device.CreatedAt,
		device.UpdatedAt,
	).Scan(&device.ID)
}

func (r *postgresDeviceRepository) Update(ctx context.Context, device *domain.Device) error {

	query := `
			UPDATE devices
			SET name = $2, token = $3, active = $4, updated_at = $5
			WHERE id = $1`

	_, err := r.conn.Exec(
		ctx,
		query,
		device.ID,
		device.Name,
		device.Token,
		device.Active,
		device.UpdatedAt,
	)
	return err
}

func (r *postgresDeviceRepository) Delete(ctx context.Context, userID, deviceID int64) error {

	query := `DELETE FROM devices WHERE id = $1 AND user_id = $2`

	_, err := r.conn.Exec(ctx, query, deviceID, userID)
	return err
}

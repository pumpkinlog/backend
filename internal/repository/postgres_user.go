package repository

import (
	"context"
	"errors"

	"github.com/pumpkinlog/backend/internal/domain"
)

type postgresUserRepository struct {
	conn Connection
}

func NewPostgresUserRepository(conn Connection) domain.UserRepository {
	return &postgresUserRepository{conn}
}

func (p *postgresUserRepository) fetch(ctx context.Context, query string, args ...any) ([]domain.User, error) {

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.FavoriteRegions,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {

	query := `
			SELECT id, favorite_regions, created_at, updated_at
			FROM users
			WHERE id = $1`

	users, err := r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, domain.ErrNotFound
	}

	return &users[0], nil
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {

	if user == nil {
		return errors.New("user is nil")
	}

	query := `
			INSERT INTO users (id, favorite_regions, created_at, updated_at)
			VALUES ($1, $2, $3, $4)`

	_, err := r.conn.Exec(
		ctx,
		query,
		user.ID,
		user.FavoriteRegions,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {

	query := `
			UPDATE users
			SET favorite_regions = $2, updated_at = $3
			WHERE id = $1`

	_, err := r.conn.Exec(
		ctx,
		query,
		user.ID,
		user.FavoriteRegions,
		user.UpdatedAt,
	)
	return err
}

func (r *postgresUserRepository) Delete(ctx context.Context, id string) error {

	query := `
			DELETE FROM users
			WHERE id = $1`

	_, err := r.conn.Exec(ctx, query, id)
	return err
}

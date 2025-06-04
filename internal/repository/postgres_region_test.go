package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/test"
)

func TestGetByID(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(ctx context.Context, t *testing.T, repo domain.RegionRepository)
		id          string
		expected    *domain.Region
		expectedErr error
	}{
		{
			name: "found region",
			setup: func(ctx context.Context, t *testing.T, repo domain.RegionRepository) {
				r := &domain.Region{
					ID:             "JE",
					Name:           "Jersey",
					Type:           domain.RegionTypeCountry,
					Continent:      domain.ContinentEurope,
					YearStartMonth: 1,
					YearStartDay:   1,
					LatLng:         [2]float64{49.21, -2.13},
				}
				require.NoError(t, repo.CreateOrUpdate(ctx, r))
			},
			id: "JE",
			expected: &domain.Region{
				ID:             "JE",
				Name:           "Jersey",
				Type:           domain.RegionTypeCountry,
				Continent:      domain.ContinentEurope,
				YearStartMonth: 1,
				YearStartDay:   1,
				LatLng:         [2]float64{49.21, -2.13},
			},
		},
		{
			name:        "not found region",
			setup:       func(ctx context.Context, t *testing.T, repo domain.RegionRepository) {},
			id:          "ZZ",
			expected:    nil,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			conn := test.NewPgxConn(t)
			tx, err := conn.Begin(ctx)
			require.NoError(t, err)

			t.Cleanup(func() {
				require.NoError(t, tx.Rollback(context.Background()))
			})

			repo := NewPostgresRegionRepository(tx)

			tc.setup(ctx, t, repo)

			region, err := repo.GetByID(ctx, tc.id)
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, region, tc.expected)
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(ctx context.Context, t *testing.T, repo domain.RegionRepository)
		filter      *domain.RegionFilter
		expected    []*domain.Region
		expectedErr error
	}{
		{
			name: "found regions",
			setup: func(ctx context.Context, t *testing.T, repo domain.RegionRepository) {
				rgns := []*domain.Region{
					{
						ID:             "JE",
						Name:           "Jersey",
						Type:           domain.RegionTypeCountry,
						Continent:      domain.ContinentEurope,
						YearStartMonth: 1,
						YearStartDay:   1,
						LatLng:         [2]float64{49.21, -2.13},
					},
					{
						ID:             "GG",
						Name:           "Guernsey",
						Type:           domain.RegionTypeCountry,
						Continent:      domain.ContinentEurope,
						YearStartMonth: 1,
						YearStartDay:   1,
						LatLng:         [2]float64{49.46, -2.58},
					},
				}
				for _, r := range rgns {
					require.NoError(t, repo.CreateOrUpdate(ctx, r))
				}
			},
			expected: []*domain.Region{
				{
					ID:             "JE",
					Name:           "Jersey",
					Type:           domain.RegionTypeCountry,
					Continent:      domain.ContinentEurope,
					YearStartMonth: 1,
					YearStartDay:   1,
					LatLng:         [2]float64{49.21, -2.13},
				},
				{
					ID:             "GG",
					Name:           "Guernsey",
					Type:           domain.RegionTypeCountry,
					Continent:      domain.ContinentEurope,
					YearStartMonth: 1,
					YearStartDay:   1,
					LatLng:         [2]float64{49.46, -2.58},
				},
			},
		},
		{
			name:  "no regions found",
			setup: func(ctx context.Context, t *testing.T, repo domain.RegionRepository) {},
			filter: &domain.RegionFilter{
				RegionIDs: []string{"ZZ"},
			},
			expected: []*domain.Region{},
		},
		{
			name: "pagination success",
			setup: func(ctx context.Context, t *testing.T, repo domain.RegionRepository) {
				rgns := []*domain.Region{
					{
						ID:             "JE",
						Name:           "Jersey",
						Type:           domain.RegionTypeCountry,
						Continent:      domain.ContinentEurope,
						YearStartMonth: 1,
						YearStartDay:   1,
						LatLng:         [2]float64{49.21, -2.13},
					},
					{
						ID:             "GG",
						Name:           "Guernsey",
						Type:           domain.RegionTypeCountry,
						Continent:      domain.ContinentEurope,
						YearStartMonth: 1,
						YearStartDay:   1,
						LatLng:         [2]float64{49.46, -2.58},
					},
				}
				for _, r := range rgns {
					require.NoError(t, repo.CreateOrUpdate(ctx, r))
				}
			},
			filter: &domain.RegionFilter{
				Page:  test.Ptr(1),
				Limit: test.Ptr(1),
			},
			expected: []*domain.Region{
				{
					ID:             "GG",
					Name:           "Guernsey",
					Type:           domain.RegionTypeCountry,
					Continent:      domain.ContinentEurope,
					YearStartMonth: 1,
					YearStartDay:   1,
					LatLng:         [2]float64{49.46, -2.58},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			conn := test.NewPgxConn(t)

			tx, err := conn.Begin(ctx)
			require.NoError(t, err)

			t.Cleanup(func() {
				require.NoError(t, tx.Rollback(context.Background()))
			})

			repo := NewPostgresRegionRepository(tx)

			tc.setup(ctx, t, repo)

			regions, err := repo.List(ctx, tc.filter)
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, regions, tc.expected)
		})
	}
}

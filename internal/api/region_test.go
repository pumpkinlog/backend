package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/test/mocks"
)

func TestGetRegion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		regionID       string
		mockGetByID    func(ctx context.Context, regionID string) (*domain.Region, error)
		expectedCode   int
		expectedRegion domain.Region
	}{
		{
			name:     "region found",
			regionID: testRegionID,
			mockGetByID: func(ctx context.Context, regionID string) (*domain.Region, error) {
				return &domain.Region{ID: regionID}, nil
			},
			expectedCode:   http.StatusOK,
			expectedRegion: domain.Region{ID: testRegionID},
		},
		{
			name:     "region not found",
			regionID: testRegionID,
			mockGetByID: func(ctx context.Context, id string) (*domain.Region, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing regionID",
			expectedCode: http.StatusNotFound,
		},
		{
			name:     "repo returns error",
			regionID: testRegionID,
			mockGetByID: func(ctx context.Context, id string) (*domain.Region, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				regionRepo: &mocks.RegionRepo{GetByIDFunc: tc.mockGetByID},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/region/%s", tc.regionID)
			req := newTestRequest(t, http.MethodGet, uri, "", "")
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got domain.Region
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedRegion, "response type incorrect")
			}
		})
	}

}

func TestListRegions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		query           url.Values
		mockList        func(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error)
		expectedCode    int
		expectedRegions []domain.Region
	}{
		{
			name: "listed regions",
			query: url.Values{
				"limit": []string{"10"},
				"page":  []string{"1"},
			},
			mockList: func(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error) {
				require.NotNil(t, filter.Limit)
				require.Equal(t, *filter.Limit, 10)
				require.NotNil(t, filter.Page)
				require.Equal(t, *filter.Page, 1)
				return make([]*domain.Region, 0), nil
			},
			expectedCode:    http.StatusOK,
			expectedRegions: make([]domain.Region, 0),
		},
		{
			name: "invalid page param type",
			query: url.Values{
				"page": []string{"invalid"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid limit param type",
			query: url.Values{
				"limit": []string{"invalid"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "out of bounds limit param",
			query: url.Values{
				"limit": []string{fmt.Sprintf("%d", PaginationMaxLimit+1)},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "default page param",
			query: url.Values{
				"limit": []string{"10"},
				"page":  []string{"0"},
			},
			mockList: func(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error) {
				require.NotNil(t, filter.Limit)
				require.Equal(t, *filter.Limit, 10)
				require.NotNil(t, filter.Page)
				require.Equal(t, *filter.Page, 1)
				return make([]*domain.Region, 0), nil
			},
			expectedCode:    http.StatusOK,
			expectedRegions: make([]domain.Region, 0),
		},
		{
			name: "repo returns error",
			mockList: func(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				regionRepo: &mocks.RegionRepo{ListFunc: tc.mockList},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/region?%s", tc.query.Encode())
			req := newTestRequest(t, http.MethodGet, uri, "", "")
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got []domain.Region
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedRegions, "response type incorrect")
			}
		})
	}
}

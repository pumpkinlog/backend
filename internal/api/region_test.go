package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/test/mocks"
)

func TestGetRegion(t *testing.T) {
	tests := []struct {
		name           string
		regionID       string
		mockGetByID    func(ctx context.Context, id string) (*domain.Region, error)
		expectedCode   int
		expectedRegion *domain.Region
	}{
		{
			name:     "region found",
			regionID: "JE",
			mockGetByID: func(ctx context.Context, id string) (*domain.Region, error) {
				return &domain.Region{ID: "JE", Name: "Jersey"}, nil
			},
			expectedCode:   http.StatusOK,
			expectedRegion: &domain.Region{ID: "JE", Name: "Jersey"},
		},
		{
			name:     "region not found",
			regionID: "JE",
			mockGetByID: func(ctx context.Context, id string) (*domain.Region, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:     "invalid region ID",
			regionID: "invalid-id",
			mockGetByID: func(ctx context.Context, id string) (*domain.Region, error) {
				panic("should not be called")
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:     "repo returns error",
			regionID: "JE",
			mockGetByID: func(ctx context.Context, id string) (*domain.Region, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &API{
				logger:     slog.New(slog.DiscardHandler),
				router:     http.NewServeMux(),
				regionRepo: mocks.RegionRepo{GetByIDFunc: tc.mockGetByID},
			}
			api.registerRoutes()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/region/%s", tc.regionID), nil)
			rr := httptest.NewRecorder()

			api.router.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("expected status code %d, got %d", tc.expectedCode, rr.Code)
			}

			if tc.expectedRegion != nil {
				var got domain.Region
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if got.ID != tc.expectedRegion.ID {
					t.Errorf("expected region ID %s, got %s", tc.expectedRegion.ID, got.ID)
				}
			}
		})
	}

}

func TestListRegions(t *testing.T) {
	tests := []struct {
		name            string
		queryParams     url.Values
		mockList        func(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error)
		expectedCode    int
		expectedRegions []*domain.Region
	}{
		{
			name: "listed regions",
			mockList: func(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error) {
				return []*domain.Region{}, nil
			},
			expectedCode:    http.StatusOK,
			expectedRegions: []*domain.Region{},
		},
		{
			name: "repo returns error",
			mockList: func(ctx context.Context, filter *domain.RegionFilter) ([]*domain.Region, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "invalid page parameter",
			queryParams: url.Values{
				"page": []string{"invalid"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid limit parameter",
			queryParams: url.Values{
				"limit": []string{"invalid"},
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &API{
				logger:     slog.New(slog.DiscardHandler),
				regionRepo: mocks.RegionRepo{ListFunc: tc.mockList},
			}

			req := httptest.NewRequest(http.MethodGet, "/region", nil)
			req.URL.RawQuery = tc.queryParams.Encode()

			rr := httptest.NewRecorder()

			api.ListRegions(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("expected status code %d, got %d", tc.expectedCode, rr.Code)
			}

			if tc.expectedRegions != nil {
				var got []*domain.Region
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				for i, region := range got {
					if region.ID != tc.expectedRegions[i].ID {
						t.Errorf("expected region ID %s, got %s", tc.expectedRegions[i].ID, region.ID)
					}
				}
			}
		})
	}
}

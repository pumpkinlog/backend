package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/test/mocks"
)

func TestEvaluateRegion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		userID             string
		regionID           string
		mockEvaluate       func(ctx context.Context, userID string, regionID string) (*domain.RegionEvaluation, error)
		expectedCode       int
		expectedEvaluation domain.RegionEvaluation
	}{
		{
			name:     "evaluation found",
			userID:   testUserID,
			regionID: testRegionID,
			mockEvaluate: func(ctx context.Context, userID, regionID string) (*domain.RegionEvaluation, error) {
				return &domain.RegionEvaluation{}, nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:     "evaluation not found",
			userID:   testUserID,
			regionID: testRegionID,
			mockEvaluate: func(ctx context.Context, userID, regionID string) (*domain.RegionEvaluation, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing userID",
			regionID:     testRegionID,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "missing regionID",
			expectedCode: http.StatusNotFound,
		},
		{
			name:     "service returns error",
			userID:   testUserID,
			regionID: testRegionID,
			mockEvaluate: func(ctx context.Context, userID, regionID string) (*domain.RegionEvaluation, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				evaluationSvc: &mocks.EvaluationService{EvaluateRegionFunc: tc.mockEvaluate},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/evaluate/%s", tc.regionID)
			req := newTestRequest(t, http.MethodGet, uri, "", tc.userID)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got domain.RegionEvaluation
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedEvaluation, "response type incorrect")
			}
		})
	}
}

func TestEvaluateRegions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		userID              string
		mockEvaluate        func(ctx context.Context, userID string) ([]*domain.RegionEvaluation, error)
		expectedCode        int
		expectedEvaluations []domain.RegionEvaluation
	}{
		{
			name:   "evaluation found",
			userID: testUserID,
			mockEvaluate: func(ctx context.Context, userID string) ([]*domain.RegionEvaluation, error) {
				return make([]*domain.RegionEvaluation, 0), nil
			},
			expectedCode:        http.StatusOK,
			expectedEvaluations: make([]domain.RegionEvaluation, 0),
		},
		{
			name:   "evaluation not found",
			userID: testUserID,
			mockEvaluate: func(ctx context.Context, userID string) ([]*domain.RegionEvaluation, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},

		{
			name:   "service returns error",
			userID: testUserID,
			mockEvaluate: func(ctx context.Context, userID string) ([]*domain.RegionEvaluation, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				evaluationSvc: &mocks.EvaluationService{EvaluateRegionsFunc: tc.mockEvaluate},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodGet, "/evaluate", "", tc.userID)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got []domain.RegionEvaluation
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedEvaluations, "response type incorrect")
			}
		})
	}
}

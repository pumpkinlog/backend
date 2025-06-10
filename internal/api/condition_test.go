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

func TestGetCondition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		mockGetByID       func(ctx context.Context, condID domain.Code) (*domain.Condition, error)
		expectedCode      int
		expectedCondition domain.Condition
	}{
		{
			name: "condition found",
			mockGetByID: func(ctx context.Context, condID domain.Code) (*domain.Condition, error) {
				return &domain.Condition{ID: condID}, nil
			},
			expectedCode:      http.StatusOK,
			expectedCondition: domain.Condition{ID: testConditionID},
		},
		{
			name: "condition not found",
			mockGetByID: func(ctx context.Context, condID domain.Code) (*domain.Condition, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "repo returns error",
			mockGetByID: func(ctx context.Context, condID domain.Code) (*domain.Condition, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				conditionSvc: &mocks.ConditionService{GetByIDFunc: tc.mockGetByID},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/condition/%s", testConditionID)
			req := newTestRequest(t, http.MethodGet, uri, "", false)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got domain.Condition
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, tc.expectedCondition, got)
			}
		})
	}
}

func TestListConditions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		query              url.Values
		mockList           func(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error)
		expectedCode       int
		expectedConditions []domain.Condition
	}{
		{
			name: "conditions found",
			query: url.Values{
				"limit": []string{"10"},
				"page":  []string{"1"},
			},
			mockList: func(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error) {
				return make([]*domain.Condition, 0), nil
			},
			expectedCode:       http.StatusOK,
			expectedConditions: make([]domain.Condition, 0),
		},
		{
			name: "repo returns error",
			mockList: func(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				conditionSvc: &mocks.ConditionService{ListFunc: tc.mockList},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/condition?%s", tc.query.Encode())
			req := newTestRequest(t, http.MethodGet, uri, "", false)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got []domain.Condition
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, tc.expectedConditions, got, "response type incorrect")
			}
		})
	}
}

package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/test/mocks"
)

func TestGetRule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockGetByID  func(ctx context.Context, ruleID int64) (*domain.Rule, error)
		expectedCode int
		expectedRule domain.Rule
	}{
		{
			name: "rule found",
			mockGetByID: func(ctx context.Context, ruleID int64) (*domain.Rule, error) {
				return &domain.Rule{ID: ruleID}, nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "rule not found",
			mockGetByID: func(ctx context.Context, ruleID int64) (*domain.Rule, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "repo returns error",
			mockGetByID: func(ctx context.Context, ruleID int64) (*domain.Rule, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				ruleRepo: &mocks.RuleRepo{GetByIDFunc: tc.mockGetByID},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodGet, "/rule/0", "", false)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got domain.Rule
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedRule, "response type incorrect")
			}
		})
	}
}

func TestListRules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockList      func(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error)
		expectedCode  int
		expectedRules []domain.Rule
	}{
		{
			name: "listed rules",
			mockList: func(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error) {
				return []*domain.Rule{}, nil
			},
			expectedCode:  http.StatusOK,
			expectedRules: make([]domain.Rule, 0),
		},
		{
			name: "repo returns error",
			mockList: func(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				ruleRepo: &mocks.RuleRepo{ListFunc: tc.mockList},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodGet, "/rule", "", false)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got []domain.Rule
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedRules, "response type incorrect")
			}
		})
	}
}

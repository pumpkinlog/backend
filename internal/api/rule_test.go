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

func TestGetRule(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockGetByID  func(ctx context.Context, ruleID domain.Code) (*domain.Rule, error)
		expectedCode int
		expectedRule domain.Rule
	}{
		{
			name: "rule found",
			mockGetByID: func(ctx context.Context, ruleID domain.Code) (*domain.Rule, error) {
				return &domain.Rule{ID: ruleID}, nil
			},
			expectedCode: http.StatusOK,
			expectedRule: domain.Rule{ID: testRuleID},
		},
		{
			name: "rule not found",
			mockGetByID: func(ctx context.Context, ruleID domain.Code) (*domain.Rule, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "repo returns error",
			mockGetByID: func(ctx context.Context, ruleID domain.Code) (*domain.Rule, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				ruleSvc: &mocks.RuleService{GetByIDFunc: tc.mockGetByID},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/rule/%s", testRuleID)
			req := newTestRequest(t, http.MethodGet, uri, "", false)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got domain.Rule
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got.ID, tc.expectedRule.ID)
				require.Equal(t, got.Name, tc.expectedRule.Name)
				require.Equal(t, got.Description, tc.expectedRule.Description)
				require.Equal(t, got.RegionID, tc.expectedRule.RegionID)
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
				ruleSvc: &mocks.RuleService{ListFunc: tc.mockList},
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

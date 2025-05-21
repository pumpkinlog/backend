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

func TestGetCondition(t *testing.T) {
	tests := []struct {
		name         string
		conditionID  string
		mockGetByID  func(ctx context.Context, id string) (*domain.Condition, error)
		expectedCode int
		expectedCond *domain.Condition
	}{
		{
			name:        "condition found",
			conditionID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e",
			mockGetByID: func(ctx context.Context, id string) (*domain.Condition, error) {
				return &domain.Condition{ID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e"}, nil
			},
			expectedCode: http.StatusOK,
			expectedCond: &domain.Condition{ID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e"},
		},
		{
			name:        "condition not found",
			conditionID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e",
			mockGetByID: func(ctx context.Context, id string) (*domain.Condition, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing condition ID",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid condition ID",
			conditionID:  "invalid-id",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "repo returns error",
			conditionID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e",
			mockGetByID: func(ctx context.Context, id string) (*domain.Condition, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &API{
				logger:        slog.New(slog.DiscardHandler),
				router:        http.NewServeMux(),
				conditionRepo: mocks.ConditionRepo{GetByIDFunc: tc.mockGetByID},
			}
			api.registerRoutes()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/condition/%s", tc.conditionID), nil)
			rr := httptest.NewRecorder()

			api.router.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("expected status code %d, got %d", tc.expectedCode, rr.Code)
			}

			if tc.expectedCond != nil {
				var got domain.Condition
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if got.ID != tc.expectedCond.ID {
					t.Errorf("expected condition ID %s, got %s", tc.expectedCond.ID, got.ID)
				}
			}
		})
	}
}

func TestListConditions(t *testing.T) {
	tests := []struct {
		name          string
		queryParams   url.Values
		mockList      func(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error)
		expectedCode  int
		expectedConds []*domain.Condition
	}{
		{
			name: "listed conditions",
			queryParams: url.Values{
				"ruleId": []string{"bf03f5b7-3a58-4555-92ae-f0bf4000052e"},
			},
			mockList: func(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error) {
				return []*domain.Condition{
					{RuleID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e"},
				}, nil
			},
			expectedCode: http.StatusOK,
			expectedConds: []*domain.Condition{
				{RuleID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e"},
			},
		},
		{
			name: "invalid rule IDs",
			queryParams: url.Values{
				"ruleId": []string{"invalid-id", "another-invalid-id"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "repo returns error",
			queryParams: url.Values{
				"ruleId": []string{"bf03f5b7-3a58-4555-92ae-f0bf4000052e"},
			},
			mockList: func(ctx context.Context, filter *domain.ConditionFilter) ([]*domain.Condition, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &API{
				logger:        slog.New(slog.DiscardHandler),
				conditionRepo: mocks.ConditionRepo{ListFunc: tc.mockList},
			}

			req := httptest.NewRequest(http.MethodGet, "/condition", nil)
			req.URL.RawQuery = tc.queryParams.Encode()

			rr := httptest.NewRecorder()

			api.ListConditions(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("expected status code %d, got %d", tc.expectedCode, rr.Code)
			}

			if tc.expectedConds != nil {
				var got []*domain.Condition
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				for i, cond := range got {
					if cond.RuleID != tc.expectedConds[i].RuleID {
						t.Errorf("expected condition RuleID %s, got %s", tc.expectedConds[i].RuleID, cond.RuleID)
					}
				}
			}
		})
	}
}

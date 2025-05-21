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

func TestGetRule(t *testing.T) {
	tests := []struct {
		name         string
		ruleID       string
		mockGetByID  func(ctx context.Context, id string) (*domain.Rule, error)
		expectedCode int
		expectedRule *domain.Rule
	}{
		{
			name:   "rule found",
			ruleID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e",
			mockGetByID: func(ctx context.Context, id string) (*domain.Rule, error) {
				return &domain.Rule{ID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e"}, nil
			},
			expectedCode: http.StatusOK,
			expectedRule: &domain.Rule{ID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e"},
		},
		{
			name:   "rule not found",
			ruleID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e",
			mockGetByID: func(ctx context.Context, id string) (*domain.Rule, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid rule ID",
			ruleID:       "invalid-id",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "repo returns error",
			ruleID: "bf03f5b7-3a58-4555-92ae-f0bf4000052e",
			mockGetByID: func(ctx context.Context, id string) (*domain.Rule, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &API{
				logger:   slog.New(slog.DiscardHandler),
				router:   http.NewServeMux(),
				ruleRepo: mocks.RuleRepo{GetByIDFunc: tc.mockGetByID},
			}
			api.registerRoutes()

			req := httptest.NewRequest("GET", fmt.Sprintf("/rule/%s", tc.ruleID), nil)
			req.SetPathValue("ruleId", tc.ruleID)

			rr := httptest.NewRecorder()

			api.GetRule(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("expected status code %d, got %d", tc.expectedCode, rr.Code)
			}

			if tc.expectedRule != nil {
				var got domain.Rule
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if got.ID != tc.expectedRule.ID {
					t.Errorf("expected rule ID %s, got %s", tc.expectedRule.ID, got.ID)
				}
			}
		})
	}
}

func TestListRules(t *testing.T) {
	tests := []struct {
		name          string
		queryParams   url.Values
		listRulesFunc func(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error)
		expectedCode  int
		expectedRules []*domain.Rule
	}{
		{
			name: "listed rules",
			queryParams: url.Values{
				"regionId": []string{"JE"},
			},
			listRulesFunc: func(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error) {
				return []*domain.Rule{}, nil
			},
			expectedCode:  http.StatusOK,
			expectedRules: []*domain.Rule{},
		},
		{
			name:         "missing region ID",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "repo returns error",
			queryParams: url.Values{
				"regionId": []string{"JE"},
			},
			listRulesFunc: func(ctx context.Context, filter *domain.RuleFilter) ([]*domain.Rule, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &API{
				logger:   slog.New(slog.DiscardHandler),
				ruleRepo: mocks.RuleRepo{ListFunc: tc.listRulesFunc},
			}

			req := httptest.NewRequest("GET", "/rule", nil)
			req.URL.RawQuery = tc.queryParams.Encode()

			rr := httptest.NewRecorder()

			api.ListRules(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("expected status code %d, got %d", tc.expectedCode, rr.Code)
			}

			if tc.expectedRules != nil {
				var got []*domain.Rule
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				for i, rule := range got {
					if rule.ID != tc.expectedRules[i].ID {
						t.Errorf("expected rule ID %s, got %s", tc.expectedRules[i].ID, rule.ID)
					}
				}
			}
		})
	}
}

package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/test/mocks"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		getUserFunc  func(ctx context.Context, id string) (*domain.User, error)
		expectedCode int
		expectedUser *domain.User
	}{
		{
			name:   "user found",
			userID: "bb262ce1-6ae4-422a-a1e3-cdcf36fd13aa",
			getUserFunc: func(ctx context.Context, id string) (*domain.User, error) {
				return &domain.User{ID: "bb262ce1-6ae4-422a-a1e3-cdcf36fd13aa"}, nil
			},
			expectedCode: http.StatusOK,
			expectedUser: &domain.User{ID: "bb262ce1-6ae4-422a-a1e3-cdcf36fd13aa"},
		},
		{
			name:   "user not found",
			userID: "bb262ce1-6ae4-422a-a1e3-cdcf36fd13aa",
			getUserFunc: func(ctx context.Context, id string) (*domain.User, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:   "repo returns error",
			userID: "bb262ce1-6ae4-422a-a1e3-cdcf36fd13aa",
			getUserFunc: func(ctx context.Context, id string) (*domain.User, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "no identity",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &API{
				logger:   slog.New(slog.DiscardHandler),
				router:   http.NewServeMux(),
				userRepo: mocks.UserRepo{GetByIDFunc: tc.getUserFunc},
			}
			api.registerRoutes()

			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			req.Header.Add("X-User-ID", tc.userID)

			rr := httptest.NewRecorder()

			api.router.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("expected status code %d, got %d", tc.expectedCode, rr.Code)
			}

			if tc.expectedUser != nil {
				var got domain.User
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if got.ID != tc.expectedUser.ID {
					t.Fatalf("expected user ID %s, got %s", tc.expectedUser.ID, got.ID)
				}
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		createUserFunc func(ctx context.Context, user *domain.User) error
		expectedCode   int
	}{
		{
			name:         "user created",
			expectedCode: http.StatusCreated,
		},
		{
			name:         "repo returns error",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Skip()
		})
	}
}

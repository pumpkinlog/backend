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

func TestGetUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		mockGetByID  func(ctx context.Context, userID string) (*domain.User, error)
		expectedCode int
		expectedUser domain.User
	}{
		{
			name:   "user found",
			userID: testUserID,
			mockGetByID: func(ctx context.Context, userID string) (*domain.User, error) {
				return &domain.User{ID: testUserID}, nil
			},
			expectedCode: http.StatusOK,
			expectedUser: domain.User{ID: testUserID},
		},
		{
			name:   "user not found",
			userID: testUserID,
			mockGetByID: func(ctx context.Context, userID string) (*domain.User, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:   "repo returns error",
			userID: testUserID,
			mockGetByID: func(ctx context.Context, userID string) (*domain.User, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				userRepo: &mocks.UserRepository{GetByIDFunc: tc.mockGetByID},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodGet, "/user", "", tc.userID)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got domain.User
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedUser, "response type incorrect")
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockCreate   func(ctx context.Context, userID string) error
		request      string
		expectedCode int
	}{
		{
			name: "user created",
			mockCreate: func(ctx context.Context, userID string) error {
				return nil
			},
			request:      fmt.Sprintf(`{"userId":"%s"}`, testUserID),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "empty body",
			request:      `{}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid body",
			request:      `{invalid json}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "empty user ID",
			request:      `{}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validation error",
			mockCreate: func(ctx context.Context, userID string) error {
				return domain.ErrValidation
			},
			request:      fmt.Sprintf(`{"userId":"%s"}`, testUserID),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "service returns error",
			mockCreate: func(ctx context.Context, userID string) error {
				return errors.New("database error")
			},
			request:      fmt.Sprintf(`{"userId":"%s"}`, testUserID),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				userSvc: &mocks.UserService{CreateFunc: tc.mockCreate},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodPost, "/user", tc.request, "")
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, rr.Body.String())
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		request      string
		mockUpdate   func(ctx context.Context, userID string, favoriteRegions, wantResidency []string) error
		expectedCode int
	}{
		{
			name:    "updated user",
			userID:  testUserID,
			request: "{}",
			mockUpdate: func(ctx context.Context, userID string, favoriteRegions, wantResidency []string) error {
				return nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:    "user not found",
			userID:  testUserID,
			request: "{}",
			mockUpdate: func(ctx context.Context, userID string, favoriteRegions, wantResidency []string) error {
				return domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid body",
			userID:       testUserID,
			request:      `{invalid json}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing user ID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:    "validation error",
			userID:  testUserID,
			request: "{}",
			mockUpdate: func(ctx context.Context, userID string, favoriteRegions, wantResidency []string) error {
				return domain.ErrValidation
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "service returns error",
			userID:  testUserID,
			request: "{}",
			mockUpdate: func(ctx context.Context, userID string, favoriteRegions, wantResidency []string) error {
				return errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				userSvc: &mocks.UserService{UpdateFunc: tc.mockUpdate},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodPatch, "/user", tc.request, tc.userID)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")
		})
	}
}

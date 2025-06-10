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

func TestGetUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		authenticated bool
		mockGetByID   func(ctx context.Context, userID int64) (*domain.User, error)
		expectedCode  int
		expectedUser  domain.User
	}{
		{
			name:          "user found",
			authenticated: true,
			mockGetByID: func(ctx context.Context, userID int64) (*domain.User, error) {
				return &domain.User{ID: userID}, nil
			},
			expectedCode: http.StatusOK,
			expectedUser: domain.User{},
		},
		{
			name:          "user not found",
			authenticated: true,
			mockGetByID: func(ctx context.Context, userID int64) (*domain.User, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "repo returns error",
			authenticated: true,
			mockGetByID: func(ctx context.Context, userID int64) (*domain.User, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				userSvc: &mocks.UserService{GetByIDFunc: tc.mockGetByID},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodGet, "/user", "", tc.authenticated)
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
		mockCreate   func(ctx context.Context, favoriteRegions, wantResidency []domain.RegionID) error
		request      string
		expectedCode int
	}{
		{
			name: "user created",
			mockCreate: func(ctx context.Context, favoriteRegions, wantResidency []domain.RegionID) error {
				return nil
			},
			request:      `{"favoriteRegions":["JE"],"wantResidency":["GG"]}`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid body",
			request:      `{invalid json}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "service returns error",
			mockCreate: func(ctx context.Context, favoriteRegions, wantResidency []domain.RegionID) error {
				return errors.New("database error")
			},
			request:      "{}",
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
			req := newTestRequest(t, http.MethodPost, "/user", tc.request, false)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, rr.Body.String())
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		authenticated bool
		request       string
		mockUpdate    func(ctx context.Context, userID int64, favoriteRegions, wantResidency []domain.RegionID) error
		expectedCode  int
	}{
		{
			name:          "updated user",
			authenticated: true,
			request:       "{}",
			mockUpdate: func(ctx context.Context, userID int64, favoriteRegions, wantResidency []domain.RegionID) error {
				return nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:          "user not found",
			authenticated: true,
			request:       "{}",
			mockUpdate: func(ctx context.Context, userID int64, favoriteRegions, wantResidency []domain.RegionID) error {
				return domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:          "invalid body",
			authenticated: true,
			request:       `{invalid json}`,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:         "missing user ID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "validation error",
			authenticated: true,
			request:       "{}",
			mockUpdate: func(ctx context.Context, userID int64, favoriteRegions, wantResidency []domain.RegionID) error {
				return domain.ErrValidation
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "service returns error",
			authenticated: true,
			request:       "{}",
			mockUpdate: func(ctx context.Context, userID int64, favoriteRegions, wantResidency []domain.RegionID) error {
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
			req := newTestRequest(t, http.MethodPatch, "/user", tc.request, tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")
		})
	}
}

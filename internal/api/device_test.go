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

func TestGetDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		authenticated  bool
		deviceID       bool
		mockGetByID    func(ctx context.Context, userID, deviceID int64) (*domain.Device, error)
		expectedCode   int
		expectedDevice domain.Device
	}{
		{
			name:          "device found",
			authenticated: true,
			deviceID:      true,
			mockGetByID: func(ctx context.Context, userID, deviceID int64) (*domain.Device, error) {
				return &domain.Device{ID: deviceID, UserID: userID}, nil
			},
			expectedCode:   http.StatusOK,
			expectedDevice: domain.Device{},
		},
		{
			name:          "device not found",
			authenticated: true,
			deviceID:      true,
			mockGetByID: func(ctx context.Context, userID, deviceID int64) (*domain.Device, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing userID",
			deviceID:     true,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "repo returns error",
			authenticated: true,
			deviceID:      true,
			mockGetByID: func(ctx context.Context, userID, deviceID int64) (*domain.Device, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				deviceSvc: &mocks.DeviceService{GetByIDFunc: tc.mockGetByID},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodGet, "/device/0", "", tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code", rr.Body.String())

			if rr.Code == http.StatusOK {
				var got domain.Device
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedDevice, "response type incorrect")
			}
		})
	}
}

func TestListDevices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		authenticated   bool
		mockList        func(ctx context.Context, userID int64) ([]*domain.Device, error)
		expectedCode    int
		expectedDevices []domain.Device
	}{
		{
			name:          "devices found",
			authenticated: true,
			mockList: func(ctx context.Context, userID int64) ([]*domain.Device, error) {
				return make([]*domain.Device, 0), nil
			},
			expectedCode:    http.StatusOK,
			expectedDevices: make([]domain.Device, 0),
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "repo returns error",
			authenticated: true,
			mockList: func(ctx context.Context, userID int64) ([]*domain.Device, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				deviceSvc: &mocks.DeviceService{ListFunc: tc.mockList},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodGet, "/device", "", tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got []domain.Device
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedDevices, "response type incorrect")
			}
		})
	}
}

func TestCreateDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		authenticated bool
		request       string
		mockCreate    func(ctx context.Context, userID int64, name, platform, model string) error
		expectedCode  int
	}{
		{
			name:          "device created",
			authenticated: true,
			request:       "{}",
			mockCreate: func(ctx context.Context, userID int64, name, platform, model string) error {
				return nil
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "invalid body",
			authenticated: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "validation error",
			authenticated: true,
			request:       `{"deviceId":0}`,
			mockCreate: func(ctx context.Context, userID int64, name, platform, model string) error {
				return domain.ErrValidation
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "service error",
			authenticated: true,
			request:       `{"deviceId":0}`,
			mockCreate: func(ctx context.Context, userID int64, name, platform, model string) error {
				return errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				deviceSvc: &mocks.DeviceService{CreateFunc: tc.mockCreate},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodPost, "/device", tc.request, tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, rr.Body.String())
		})
	}
}

func TestUpdateDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		authenticated bool
		request       string
		mockUpdate    func(ctx context.Context, userID, deviceID int64, name, token string, acive bool) error
		expectedCode  int
	}{
		{
			name:          "updated device",
			authenticated: true,
			request:       "{}",
			mockUpdate: func(ctx context.Context, userID, deviceID int64, name, token string, acive bool) error {
				return nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:          "device not found",
			authenticated: true,
			request:       `{"deviceId":0}`,
			mockUpdate: func(ctx context.Context, userID, deviceID int64, name, token string, acive bool) error {
				return domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "invalid body",
			authenticated: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "validation error",
			authenticated: true,
			request:       "{}",
			mockUpdate: func(ctx context.Context, userID, deviceID int64, name, token string, acive bool) error {
				return domain.ErrValidation
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "service returns error",
			authenticated: true,
			request:       `{"deviceId":0}`,
			mockUpdate: func(ctx context.Context, userID, deviceID int64, name, token string, acive bool) error {
				return errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				deviceSvc: &mocks.DeviceService{UpdateFunc: tc.mockUpdate},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodPatch, "/device", tc.request, tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, rr.Body.String())
		})
	}
}

func TestDeleteDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		authenticated bool
		mockDelete    func(ctx context.Context, userID, deviceID int64) error
		expectedCode  int
	}{
		{
			name:          "device deleted",
			authenticated: true,
			mockDelete: func(ctx context.Context, userID, deviceID int64) error {
				return nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "repo returns error",
			authenticated: true,
			mockDelete: func(ctx context.Context, userID, deviceID int64) error {
				return errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				deviceSvc: &mocks.DeviceService{DeleteFunc: tc.mockDelete},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodDelete, "/device/0", "", tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code", rr.Body.String())
		})
	}
}

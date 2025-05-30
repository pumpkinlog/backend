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

func TestGetDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		userID         string
		deviceID       string
		mockGetByID    func(ctx context.Context, userID, deviceID string) (*domain.Device, error)
		expectedCode   int
		expectedDevice domain.Device
	}{
		{
			name:     "device found",
			userID:   testUserID,
			deviceID: testDeviceID,
			mockGetByID: func(ctx context.Context, userID, deviceID string) (*domain.Device, error) {
				return &domain.Device{ID: deviceID, UserID: userID}, nil
			},
			expectedCode:   http.StatusOK,
			expectedDevice: domain.Device{ID: testDeviceID, UserID: testUserID},
		},
		{
			name:     "device not found",
			userID:   testUserID,
			deviceID: testDeviceID,
			mockGetByID: func(ctx context.Context, userID, deviceID string) (*domain.Device, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing userID",
			deviceID:     testDeviceID,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "missing deviceID",
			expectedCode: http.StatusNotFound,
		},
		{
			name:     "repo returns error",
			userID:   testUserID,
			deviceID: testDeviceID,
			mockGetByID: func(ctx context.Context, userID, deviceID string) (*domain.Device, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				deviceRepo: &mocks.DeviceRepository{GetByIDFunc: tc.mockGetByID},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/device/%s", tc.deviceID)
			req := newTestRequest(t, http.MethodGet, uri, "", tc.userID)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

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
		userID          string
		mockList        func(ctx context.Context, userID string) ([]*domain.Device, error)
		expectedCode    int
		expectedDevices []domain.Device
	}{
		{
			name:   "devices found",
			userID: testUserID,
			mockList: func(ctx context.Context, userID string) ([]*domain.Device, error) {
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
			name:   "repo returns error",
			userID: testUserID,
			mockList: func(ctx context.Context, userID string) ([]*domain.Device, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				deviceRepo: &mocks.DeviceRepository{ListFunc: tc.mockList},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodGet, "/device", "", tc.userID)
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
		name         string
		userID       string
		request      string
		mockCreate   func(ctx context.Context, userID, name, platform, model string) error
		expectedCode int
	}{
		{
			name:    "device created",
			userID:  testUserID,
			request: "{}",
			mockCreate: func(ctx context.Context, userID, name, platform, model string) error {
				return nil
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid body",
			userID:       testUserID,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "validation error",
			userID:  testUserID,
			request: fmt.Sprintf(`{"deviceId":"%s"}`, testDeviceID),
			mockCreate: func(ctx context.Context, userID, name, platform, model string) error {
				return domain.ErrValidation
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "service error",
			userID:  testUserID,
			request: fmt.Sprintf(`{"deviceId":"%s"}`, testDeviceID),
			mockCreate: func(ctx context.Context, userID, name, platform, model string) error {
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
			req := newTestRequest(t, http.MethodPost, "/device", tc.request, tc.userID)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, rr.Body.String())
		})
	}
}

func TestUpdateDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		request      string
		mockUpdate   func(ctx context.Context, userID, deviceID, name, token string, acive bool) error
		expectedCode int
	}{
		{
			name:    "updated device",
			userID:  testUserID,
			request: "{}",
			mockUpdate: func(ctx context.Context, userID, deviceID, name, token string, acive bool) error {
				return nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:    "device not found",
			userID:  testUserID,
			request: fmt.Sprintf(`{"deviceId":"%s"}`, testDeviceID),
			mockUpdate: func(ctx context.Context, userID, deviceID, name, token string, acive bool) error {
				return domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid body",
			userID:       testUserID,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "validation error",
			userID:  testUserID,
			request: "{}",
			mockUpdate: func(ctx context.Context, userID, deviceID, name, token string, acive bool) error {
				return domain.ErrValidation
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "service returns error",
			userID:  testUserID,
			request: fmt.Sprintf(`{"deviceId":"%s"}`, testDeviceID),
			mockUpdate: func(ctx context.Context, userID, deviceID, name, token string, acive bool) error {
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
			req := newTestRequest(t, http.MethodPatch, "/device", tc.request, tc.userID)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, rr.Body.String())
		})
	}
}

func TestDeleteDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		deviceID     string
		mockDelete   func(ctx context.Context, userID, deviceID string) error
		expectedCode int
	}{
		{
			name:     "device deleted",
			userID:   testUserID,
			deviceID: testDeviceID,
			mockDelete: func(ctx context.Context, userID, deviceID string) error {
				return nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "missing userID",
			deviceID:     testDeviceID,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "missing deviceID",
			expectedCode: http.StatusNotFound,
		},
		{
			name:     "repo returns error",
			userID:   testUserID,
			deviceID: testDeviceID,
			mockDelete: func(ctx context.Context, userID, deviceID string) error {
				return errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				deviceRepo: &mocks.DeviceRepository{DeleteFunc: tc.mockDelete},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/device/%s", tc.deviceID)
			req := newTestRequest(t, http.MethodDelete, uri, "", tc.userID)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")
		})
	}
}

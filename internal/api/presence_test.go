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
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/test/mocks"
)

func TestGetPresence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		authenticated    bool
		regionID         string
		date             string
		mockGetByID      func(ctx context.Context, userID int64, regionID string, date time.Time) (*domain.Presence, error)
		expectedCode     int
		expectedPresence domain.Presence
	}{
		{
			name:          "presence found",
			authenticated: true,
			regionID:      testRegionID,
			date:          testDate.Format(time.DateOnly),
			mockGetByID: func(ctx context.Context, userID int64, regionID string, date time.Time) (*domain.Presence, error) {
				return &domain.Presence{
					UserID:   userID,
					RegionID: regionID,
					Date:     date,
				}, nil
			},
			expectedCode: http.StatusOK,
			expectedPresence: domain.Presence{
				RegionID: testRegionID,
				Date:     testDate,
			},
		},
		{
			name:          "presence not found",
			authenticated: true,
			regionID:      testRegionID,
			date:          testDate.Format(time.DateOnly),
			mockGetByID: func(ctx context.Context, userID int64, regionID string, date time.Time) (*domain.Presence, error) {
				return nil, domain.ErrNotFound
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "missing userID",
			regionID:     testRegionID,
			date:         testDate.Format(time.DateOnly),
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "invalid date",
			authenticated: true,
			date:          "2025-01-35",
			regionID:      testRegionID,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "repo returns error",
			authenticated: true,
			regionID:      testRegionID,
			date:          testDate.Format(time.DateOnly),
			mockGetByID: func(ctx context.Context, userID int64, regionID string, date time.Time) (*domain.Presence, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				presenceRepo: &mocks.PresenceRepo{GetByIDFunc: tc.mockGetByID},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/presence/%s/%s", tc.regionID, tc.date)
			req := newTestRequest(t, http.MethodGet, uri, "", tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got domain.Presence
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedPresence, "response type incorrect")
			}
		})
	}
}

func TestListPresence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		authenticated     bool
		query             url.Values
		mockList          func(ctx context.Context, userID int64, filter *domain.PresenceFilter) ([]*domain.Presence, error)
		expectedCode      int
		expectedPresences []domain.Presence
	}{
		{
			name:          "listed presences",
			authenticated: true,
			query: url.Values{
				"start": []string{testDate.Format(time.DateOnly)},
				"end":   []string{testDate.Format(time.DateOnly)},
				"limit": []string{"10"},
				"page":  []string{"1"},
			},
			mockList: func(ctx context.Context, userID int64, filter *domain.PresenceFilter) ([]*domain.Presence, error) {
				require.NotNil(t, filter.Limit)
				require.Equal(t, *filter.Limit, 10)
				require.NotNil(t, filter.Page)
				require.Equal(t, *filter.Page, 1)
				require.NotNil(t, filter.Start)
				require.Equal(t, *filter.Start, testDate)
				require.NotNil(t, filter.End)
				require.Equal(t, *filter.End, testDate)
				return make([]*domain.Presence, 0), nil
			},
			expectedCode:      http.StatusOK,
			expectedPresences: make([]domain.Presence, 0),
		},
		{
			name:          "invalid start param",
			authenticated: true,
			query: url.Values{
				"start": []string{"invalid time"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "invalid end param",
			authenticated: true,
			query: url.Values{
				"start": []string{testDate.Format(time.DateOnly)},
				"end":   []string{"invalid time"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "invalid page param type",
			authenticated: true,
			query: url.Values{
				"page": []string{"invalid"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "invalid limit param type",
			authenticated: true,
			query: url.Values{
				"limit": []string{"invalid"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "out of bounds limit param",
			authenticated: true,
			query: url.Values{
				"limit": []string{fmt.Sprintf("%d", PaginationMaxLimit+1)},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "default page param",
			authenticated: true,
			query: url.Values{
				"limit": []string{"10"},
				"page":  []string{"0"},
			},
			mockList: func(ctx context.Context, userID int64, filter *domain.PresenceFilter) ([]*domain.Presence, error) {
				require.NotNil(t, filter.Limit)
				require.Equal(t, *filter.Limit, 10)
				require.NotNil(t, filter.Page)
				require.Equal(t, *filter.Page, 1)
				return make([]*domain.Presence, 0), nil
			},
			expectedCode:      http.StatusOK,
			expectedPresences: make([]domain.Presence, 0),
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "repo returns error",
			authenticated: true,
			mockList: func(ctx context.Context, userID int64, filter *domain.PresenceFilter) ([]*domain.Presence, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				presenceRepo: &mocks.PresenceRepo{ListFunc: tc.mockList},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/presence?%s", tc.query.Encode())
			req := newTestRequest(t, http.MethodGet, uri, "", tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")

			if rr.Code == http.StatusOK {
				var got []domain.Presence
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedPresences, "response type incorrect")
			}
		})
	}
}

func TestCreatePresence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		authenticated bool
		request       string
		mockCreate    func(ctx context.Context, userID int64, regionID string, deviceID *int64, start, end time.Time) error
		expectedCode  int
	}{
		{
			name:          "created presence",
			authenticated: true,
			request:       fmt.Sprintf(`{"regionId":"%s","start":"%s","end":"%s"}`, testRegionID, testDate.Format(time.DateOnly), testDate.Format(time.DateOnly)),
			mockCreate: func(ctx context.Context, userID int64, regionID string, deviceID *int64, start, end time.Time) error {
				return nil
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "missing user ID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "invalid start",
			authenticated: true,
			request:       "{}",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "invalid end",
			authenticated: true,
			request:       fmt.Sprintf(`{"start":"%s"}`, testDate.Format(time.DateOnly)),
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "validation error",
			authenticated: true,
			mockCreate: func(ctx context.Context, userID int64, regionID string, deviceID *int64, start, end time.Time) error {
				return domain.ErrValidation
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "service error",
			authenticated: true,
			request:       fmt.Sprintf(`{"regionId":"%s","start":"%s","end":"%s"}`, testRegionID, testDate.Format(time.DateOnly), testDate.Format(time.DateOnly)),
			mockCreate: func(ctx context.Context, userID int64, regionID string, deviceID *int64, start, end time.Time) error {
				return errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				presenceSvc: &mocks.PresenceService{CreateFunc: tc.mockCreate},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodPost, "/presence", tc.request, tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")
		})
	}
}

func TestDeletePresence(t *testing.T) {
	t.Parallel()

	dateStr := testDate.Format(time.DateOnly)

	tests := []struct {
		name          string
		authenticated bool
		query         url.Values
		mockDelete    func(ctx context.Context, userID int64, regionID string, start, end time.Time) error
		expectedCode  int
	}{
		{
			name:          "deleted presence",
			authenticated: true,
			query: url.Values{
				"regionId": []string{testRegionID},
				"start":    []string{dateStr},
				"end":      []string{dateStr},
			},
			mockDelete: func(ctx context.Context, userID int64, regionID string, start, end time.Time) error {
				return nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "missing userID",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:          "missing regionID",
			authenticated: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "invalid start",
			authenticated: true,
			query: url.Values{
				"regionId": []string{testRegionID},
				"start":    []string{"invalid time"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "invalid end",
			authenticated: true,
			query: url.Values{
				"regionId": []string{testRegionID},
				"start":    []string{testDate.Format(time.DateOnly)},
				"end":      []string{"invalid time"},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "repo returns error",
			authenticated: true,
			query: url.Values{
				"regionId": []string{testRegionID},
				"start":    []string{dateStr},
				"end":      []string{dateStr},
			},
			mockDelete: func(ctx context.Context, userID int64, regionID string, start, end time.Time) error {
				return errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				presenceSvc: &mocks.PresenceService{DeleteFunc: tc.mockDelete},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/presence?%s", tc.query.Encode())
			req := newTestRequest(t, http.MethodDelete, uri, "", tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code", rr.Body.String())
		})
	}
}

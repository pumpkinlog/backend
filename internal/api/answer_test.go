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

func TestGetAnswer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		authenticated  bool
		mockGetByID    func(ctx context.Context, userID int64, conditionID domain.Code) (*domain.Answer, error)
		expectedCode   int
		expectedAnswer domain.Answer
	}{
		{
			name:          "answer found",
			authenticated: true,
			mockGetByID: func(ctx context.Context, userID int64, conditionID domain.Code) (*domain.Answer, error) {
				return &domain.Answer{}, nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:          "answer not found",
			authenticated: true,
			mockGetByID: func(ctx context.Context, userID int64, conditionID domain.Code) (*domain.Answer, error) {
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
			mockGetByID: func(ctx context.Context, userID int64, conditionID domain.Code) (*domain.Answer, error) {
				return nil, errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				answerSvc: &mocks.AnswerService{GetByIDFunc: tc.mockGetByID},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/answer/%s", testConditionID)
			req := newTestRequest(t, http.MethodGet, uri, "", tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code", rr.Body.String())

			if rr.Code == http.StatusOK {
				var got domain.Answer
				err := json.NewDecoder(rr.Body).Decode(&got)
				require.NoError(t, err, "cannot decode json response")
				require.Equal(t, got, tc.expectedAnswer, "response type incorrect")
			}
		})
	}
}

func TestSubmitAnswer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		authenticated      bool
		request            string
		mockCreateOrUpdate func(ctx context.Context, userID int64, conditionID domain.Code, value any) error
		expectedCode       int
	}{
		{
			name:          "successful submission",
			authenticated: true,
			request:       "{}",
			mockCreateOrUpdate: func(ctx context.Context, userID int64, conditionID domain.Code, value any) error {
				return nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "missing user ID",
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
			mockCreateOrUpdate: func(ctx context.Context, userID int64, conditionID domain.Code, value any) error {
				return domain.ErrValidation
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "service error",
			authenticated: true,
			request:       "{}",
			mockCreateOrUpdate: func(ctx context.Context, userID int64, conditionID domain.Code, value any) error {
				return errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				answerSvc: &mocks.AnswerService{CreateOrUpdateFunc: tc.mockCreateOrUpdate},
			}

			api := newTestAPI(t, opts)
			req := newTestRequest(t, http.MethodPost, "/answer", tc.request, tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code")
		})
	}
}

func TestDeleteAnswer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		authenticated  bool
		mockDelete     func(ctx context.Context, userID int64, conditionID domain.Code) error
		expectedCode   int
		expectedAnswer domain.Answer
	}{
		{
			name:          "deleted answer",
			authenticated: true,
			mockDelete: func(ctx context.Context, userID int64, conditionID domain.Code) error {
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
			mockDelete: func(ctx context.Context, userID int64, conditionID domain.Code) error {
				return errors.New("database error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := testAPIOptions{
				answerSvc: &mocks.AnswerService{DeleteFunc: tc.mockDelete},
			}

			api := newTestAPI(t, opts)
			uri := fmt.Sprintf("/answer/%s", testConditionID)
			req := newTestRequest(t, http.MethodDelete, uri, "", tc.authenticated)
			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code, "unexpected status code", rr.Body.String())
		})
	}
}

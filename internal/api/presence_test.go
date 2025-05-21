package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/test/mocks"
)

func TestGetPresence(t *testing.T) {

	tests := []struct {
		name             string
		userID           string
		regionID         string
		date             time.Time
		getPresenceFunc  func(ctx context.Context, userID, regionID string, date time.Time) (*domain.Presence, error)
		expectedCode     int
		expectedPresence *domain.Presence
	}{
		{
			name:     "presence found",
			userID:   "df228dbd-ca3f-47d1-9256-a8fea90d95d5",
			regionID: "JE",
			date:     time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			getPresenceFunc: func(ctx context.Context, userID, regionID string, date time.Time) (*domain.Presence, error) {
				return &domain.Presence{
					UserID:   "df228dbd-ca3f-47d1-9256-a8fea90d95d5",
					RegionID: "JE",
					Date:     time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
				}, nil
			},
			expectedCode: http.StatusOK,
			expectedPresence: &domain.Presence{
				UserID:   "df228dbd-ca3f-47d1-9256-a8fea90d95d5",
				RegionID: "JE",
				Date:     time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			api := &API{
				logger:       slog.New(slog.DiscardHandler),
				router:       http.NewServeMux(),
				presenceRepo: mocks.PresenceRepo{GetByIDFunc: tc.getPresenceFunc},
			}
			api.registerRoutes()

			uri := fmt.Sprintf("/presence/%s/%s", tc.regionID, tc.date.Format(time.DateOnly))
			fmt.Println(uri)
			req := httptest.NewRequest("GET", uri, nil)
			req.Header.Add("X-User-ID", tc.userID)

			rr := httptest.NewRecorder()

			api.router.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("expected status code %d, got %d", tc.expectedCode, rr.Code)
			}

			if tc.expectedPresence != nil {
				var got domain.Presence
				if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				require.Equal(t, *tc.expectedPresence, got)
			}
		})
	}
}

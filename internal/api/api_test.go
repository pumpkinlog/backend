package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pumpkinlog/backend/internal/domain"
)

var (
	testRegionID    = domain.RegionID("JE")
	testRuleID      = domain.Code("test_rule_id")
	testConditionID = domain.Code("test_condition_id")
	testDate        = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
)

type testAPIOptions struct {
	userSvc       domain.UserService
	presenceSvc   domain.PresenceService
	deviceSvc     domain.DeviceService
	answerSvc     domain.AnswerService
	evaluationSvc domain.EvaluationService
	conditionSvc  domain.ConditionService
	regionSvc     domain.RegionService
	ruleSvc       domain.RuleService
}

func newTestAPI(t *testing.T, opts testAPIOptions) *API {
	t.Helper()

	a := &API{
		logger: slog.New(slog.DiscardHandler),
		router: http.NewServeMux(),

		userSvc:       opts.userSvc,
		regionSvc:     opts.regionSvc,
		presenceSvc:   opts.presenceSvc,
		deviceSvc:     opts.deviceSvc,
		answerSvc:     opts.answerSvc,
		conditionSvc:  opts.conditionSvc,
		ruleSvc:       opts.ruleSvc,
		evaluationSvc: opts.evaluationSvc,
	}

	a.registerRoutes()

	return a
}

func newTestRequest(t *testing.T, method, path string, body string, authenticated bool) *http.Request {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if authenticated {
		req.Header.Set("X-User-ID", "0")
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

func TestRespondJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := map[string]string{"foo": "bar"}

	RespondJSON(rr, http.StatusTeapot, payload)

	require.Equal(t, http.StatusTeapot, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var got map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &got)
	require.NoError(t, err)
	require.Equal(t, payload, got)
}

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
	testRegionID = "JE"
	testDate     = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
)

type testAPIOptions struct {
	userSvc       domain.UserService
	presenceSvc   domain.PresenceService
	deviceSvc     domain.DeviceService
	answerSvc     domain.AnswerService
	evaluationSvc domain.EvaluationService

	userRepo      domain.UserRepository
	regionRepo    domain.RegionRepository
	presenceRepo  domain.PresenceRepository
	deviceRepo    domain.DeviceRepository
	ruleRepo      domain.RuleRepository
	conditionRepo domain.ConditionRepository
	answerRepo    domain.AnswerRepository
}

func newTestAPI(t *testing.T, opts testAPIOptions) *API {
	t.Helper()

	return &API{
		logger: slog.New(slog.DiscardHandler),
		router: http.NewServeMux(),

		userSvc:       opts.userSvc,
		presenceSvc:   opts.presenceSvc,
		deviceSvc:     opts.deviceSvc,
		answerSvc:     opts.answerSvc,
		evaluationSvc: opts.evaluationSvc,

		userRepo:      opts.userRepo,
		regionRepo:    opts.regionRepo,
		presenceRepo:  opts.presenceRepo,
		deviceRepo:    opts.deviceRepo,
		ruleRepo:      opts.ruleRepo,
		conditionRepo: opts.conditionRepo,
		answerRepo:    opts.answerRepo,
	}
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

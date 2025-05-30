package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rabbitmq/amqp091-go"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/repository"
	"github.com/pumpkinlog/backend/internal/service"
)

const (
	PaginationMaxLimit     = 100
	PaginationDefaultLimit = 25
)

type API struct {
	logger *slog.Logger
	router *http.ServeMux
	ch     *amqp091.Channel

	userSvc       domain.UserService
	deviceSvc     domain.DeviceService
	answerSvc     domain.AnswerService
	presenceSvc   domain.PresenceService
	evaluationSvc domain.EvaluationService

	regionRepo    domain.RegionRepository
	ruleRepo      domain.RuleRepository
	answerRepo    domain.AnswerRepository
	conditionRepo domain.ConditionRepository
	deviceRepo    domain.DeviceRepository
	presenceRepo  domain.PresenceRepository
	userRepo      domain.UserRepository
}

func NewAPI(logger *slog.Logger, conn repository.Executor, ch *amqp091.Channel) *API {
	return &API{
		logger: logger,
		router: http.NewServeMux(),
		ch:     ch,

		userSvc:       service.NewUserService(conn),
		deviceSvc:     service.NewDeviceService(conn),
		answerSvc:     service.NewAnswerService(conn),
		presenceSvc:   service.NewPresenceService(conn, ch),
		evaluationSvc: service.NewEvaluationService(logger, conn),

		regionRepo:    repository.NewPostgresRegionRepository(conn),
		ruleRepo:      repository.NewPostgresRuleRepository(conn),
		answerRepo:    repository.NewPostgresAnswerRepository(conn),
		conditionRepo: repository.NewPostgresConditionRepository(conn),
		deviceRepo:    repository.NewPostgresDeviceRepository(conn),
		presenceRepo:  repository.NewPostgresPresenceRepository(conn),
		userRepo:      repository.NewPostgresUserRepository(conn),
	}
}

func (a *API) registerRoutes() {

	a.handle("GET /docs/openapi.yaml", a.ServeSpec, a.Logging)
	a.handle("GET /docs", a.ServeUI, a.Logging)

	a.handle("GET /region/{regionId}", a.GetRegion, a.Logging)
	a.handle("GET /region", a.ListRegions, a.Logging)

	a.handle("GET /rule/{ruleId}", a.GetRule, a.Logging)
	a.handle("GET /rule", a.ListRules, a.Logging)

	a.handle("GET /answer/{conditionId}", a.GetAnswer, a.Logging, a.Auth)
	a.handle("POST /answer", a.SubmitAnswer, a.Logging, a.Auth)
	a.handle("DELETE /answer/{conditionId}", a.DeleteAnswer, a.Logging, a.Auth)

	a.handle("GET /evaluate/{regionId}", a.EvaluateRegion, a.Logging, a.Auth)
	a.handle("GET /evaluate", a.EvaluateRegions, a.Logging, a.Auth)

	a.handle("GET /condition/{conditionId}", a.GetCondition, a.Logging)
	a.handle("GET /condition", a.ListConditions, a.Logging)

	a.handle("GET /device/{deviceId}", a.GetDevice, a.Logging, a.Auth)
	a.handle("GET /device", a.ListDevices, a.Logging, a.Auth)
	a.handle("POST /device", a.CreateDevice, a.Logging, a.Auth)
	a.handle("PATCH /device", a.UpdateDevice, a.Logging, a.Auth)
	a.handle("DELETE /device/{id}", a.DeleteDevice, a.Logging, a.Auth)

	a.handle("GET /presence/{regionId}/{date}", a.GetPresence, a.Logging, a.Auth)
	a.handle("GET /presence", a.ListPresences, a.Logging, a.Auth)
	a.handle("POST /presence", a.CreatePresence, a.Logging, a.Auth)
	a.handle("DELETE /presence", a.DeletePresence, a.Logging, a.Auth)

	a.handle("GET /user", a.GetUser, a.Logging, a.Auth)
	a.handle("POST /user", a.CreateUser, a.Logging)
	a.handle("PATCH /user", a.UpdateUser, a.Logging, a.Auth)
}

func (a *API) Handler() http.Handler {
	a.registerRoutes()
	return a.router
}

func (a *API) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Handler(),
	}
}

func (a *API) handle(pattern string, h http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	a.router.HandleFunc(pattern, h)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func RespondError(w http.ResponseWriter, code int, msg string) {
	RespondJSON(w, code, ErrorResponse{
		Message: msg,
	})
}

func RespondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type contextKey string

const UserIDKey contextKey = "userId"

func UserID(r *http.Request) string {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userID == "" {
		// can safely panic here. Will only occur if auth middleware is not set on a route before calling.
		panic("user ID not present in context")
	}
	return userID
}

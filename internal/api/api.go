package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/service"
)

type Middleware func(http.Handler) http.Handler

type API struct {
	logger     *slog.Logger
	ch         *amqp091.Channel
	router     *http.ServeMux
	middleware []Middleware

	userSvc       domain.UserService
	regionSvc     domain.RegionService
	deviceSvc     domain.DeviceService
	answerSvc     domain.AnswerService
	presenceSvc   domain.PresenceService
	evaluationSvc domain.EvaluationService
	conditionSvc  domain.ConditionService
	ruleSvc       domain.RuleService
}

func NewAPI(logger *slog.Logger, conn *pgxpool.Pool, ch *amqp091.Channel) *API {
	api := &API{
		logger: logger,
		router: http.NewServeMux(),
		ch:     ch,

		userSvc:       service.NewUserService(logger, conn),
		regionSvc:     service.NewRegionService(logger, conn),
		deviceSvc:     service.NewDeviceService(conn),
		answerSvc:     service.NewAnswerService(logger, conn),
		presenceSvc:   service.NewPresenceService(logger, conn, ch),
		conditionSvc:  service.NewConditionService(logger, conn),
		ruleSvc:       service.NewRuleService(logger, conn),
		evaluationSvc: service.NewEvaluationService(logger, conn, ch),
	}

	api.use(api.Logging, api.Cors)
	api.registerRoutes()

	return api
}

func (a *API) registerRoutes() {
	a.handle("GET /docs/openapi.yaml", a.ServeSpec)
	a.handle("GET /docs", a.ServeUI)

	a.handle("GET /region/{regionId}", a.GetRegion)
	a.handle("GET /region", a.ListRegions)

	a.handle("GET /rule/{ruleId}", a.GetRule)
	a.handle("GET /rule", a.ListRules)

	a.handle("GET /answer/{conditionId}", a.GetAnswer, a.Auth)
	a.handle("POST /answer", a.SubmitAnswer, a.Auth)
	a.handle("DELETE /answer/{conditionId}", a.DeleteAnswer, a.Auth)

	a.handle("GET /evaluate/{regionId}", a.EvaluateRegion, a.Auth)
	a.handle("GET /evaluate", a.EvaluateRegions, a.Auth)

	a.handle("GET /condition/{conditionId}", a.GetCondition)
	a.handle("GET /condition", a.ListConditions)

	a.handle("GET /device/{deviceId}", a.GetDevice, a.Auth)
	a.handle("GET /device", a.ListDevices, a.Auth)
	a.handle("POST /device", a.CreateDevice, a.Auth)
	a.handle("PATCH /device", a.UpdateDevice, a.Auth)
	a.handle("DELETE /device/{deviceId}", a.DeleteDevice, a.Auth)

	a.handle("GET /presence/{regionId}/{date}", a.GetPresence, a.Auth)
	a.handle("GET /presence", a.ListPresences, a.Auth)
	a.handle("POST /presence", a.CreatePresence, a.Auth)
	a.handle("DELETE /presence", a.DeletePresence, a.Auth)

	a.handle("GET /user", a.GetUser, a.Auth)
	a.handle("POST /user", a.CreateUser)
	a.handle("PATCH /user", a.UpdateUser, a.Auth)
}

func (a *API) Handler() http.Handler {
	return chain(a.router, a.middleware...)
}

func (a *API) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Handler(),
	}
}

func (a *API) use(mw ...Middleware) {
	a.middleware = append(a.middleware, mw...)
}

func chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func (a *API) handle(pattern string, handler http.HandlerFunc, mws ...Middleware) {
	final := chain(http.HandlerFunc(handler), mws...)
	a.router.Handle(pattern, final)
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

// UserID returns the authenticated user ID from the context injected by the auth middleware.
//
// Will panic if used on a route with no auth middleware.
func UserID(ctx context.Context) int64 {
	uid, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		// We can safely panic here as this will only occur if auth middleware is not set.
		// The app should crash to prevent unauthorized access to protected resources.
		panic("user ID is missing from context")
	}
	return uid
}

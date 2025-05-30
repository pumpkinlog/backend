package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type LoggingResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
	bytes      int
}

func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.w.Header()
}

func (lrw *LoggingResponseWriter) Write(bb []byte) (int, error) {
	wb, err := lrw.w.Write(bb)
	lrw.bytes += wb
	return wb, err
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.w.WriteHeader(statusCode)
	lrw.statusCode = statusCode
}

func (a *API) Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		lrw := &LoggingResponseWriter{w: w}

		next(lrw, r)

		duration := time.Since(start)

		key := fmt.Sprintf("%s %s", r.Method, r.RequestURI)
		metric_endpoint_invocations.Add(key, 1)

		a.logger.Info("request",
			"duration", duration.String(),
			"method", r.Method,
			"status", lrw.statusCode,
			"requestUri", r.RequestURI,
			"remoteAddr", r.RemoteAddr,
		)
	}
}

func (a *API) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// User ID should be injected by auth service
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			a.logger.Warn("authenticated route is missing user ID", "requestUri", r.RequestURI)
			RespondError(w, http.StatusUnauthorized, "missing user ID")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

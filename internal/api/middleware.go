package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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

		a.logger.Info("served request",
			"status", lrw.statusCode,
			"method", r.Method,
			"uri", r.RequestURI,
			"remoteAddr", r.RemoteAddr,
			"duration", duration.String(),
		)
	}
}

func (a *API) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// User ID should be injected by auth service
		userIDHeader := r.Header.Get("X-User-ID")
		if userIDHeader == "" {
			a.logger.Warn("authenticated route is missing user ID", "requestUri", r.RequestURI)
			RespondError(w, http.StatusUnauthorized, "missing user ID")
			return
		}

		userID, err := strconv.ParseInt(userIDHeader, 10, 64)
		if err != nil {
			a.logger.Warn("invalid user ID format", "userId", userIDHeader, "error", err)
			RespondError(w, http.StatusUnauthorized, "invalid user ID format")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServeSpec(t *testing.T) {
	t.Parallel()

	api := newTestAPI(t, testAPIOptions{})
	req := newTestRequest(t, http.MethodGet, "/docs/openapi.yaml", "", false)
	rr := httptest.NewRecorder()
	api.Handler().ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusOK, "unexpected status code")
	require.Equal(t, rr.Header().Get("Content-Type"), "text/yaml")
}

func TestServeUI(t *testing.T) {
	t.Parallel()

	api := newTestAPI(t, testAPIOptions{})
	req := newTestRequest(t, http.MethodGet, "/docs", "", false)
	rr := httptest.NewRecorder()
	api.Handler().ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusOK, "unexpected status code")
	require.Equal(t, rr.Header().Get("Content-Type"), "text/html")
}

package api

import (
	"net/http"

	"github.com/pumpkinlog/backend/docs"
)

func (a *API) ServeSpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(docs.Spec)
}

func (a *API) ServeUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(docs.UI)
}

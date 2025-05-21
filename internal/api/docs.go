package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/pumpkinlog/backend/internal/apiutil"
)

func (a *API) ServeSpec(w http.ResponseWriter, r *http.Request) {

	rootDir, err := os.Getwd()
	if err != nil {
		a.logger.Error("failed to get current directory", "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to get current directory")
		return
	}

	w.Header().Set("Content-Type", "text/yaml")

	http.ServeFile(w, r, filepath.Join(rootDir, "docs", "openapi.yml"))
}

func (a *API) ServeUI(w http.ResponseWriter, r *http.Request) {

	rootDir, err := os.Getwd()
	if err != nil {
		a.logger.Error("failed to get current directory", "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to get current directory")
		return
	}

	w.Header().Set("Content-Type", "text/html")

	http.ServeFile(w, r, filepath.Join(rootDir, "docs", "ui.html"))
}

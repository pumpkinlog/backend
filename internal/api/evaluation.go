package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) EvaluateRegion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserID(ctx)
	regionID := domain.RegionID(r.PathValue("regionId"))

	var err error
	var pit time.Time
	if pitStr := r.URL.Query().Get("pointInTime"); pitStr != "" {
		pit, err = time.Parse(time.DateOnly, pitStr)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid point in time format")
			return
		}
	}

	opts := &domain.EvaluateOpts{
		PointInTime: pit,
	}

	evaluation, err := a.evaluationSvc.EvaluateRegion(ctx, userID, regionID, opts)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondJSON(w, http.StatusNotFound, "region not found")
		default:
			a.logger.Error("failed to evaluate region", "userId", userID, "regionId", regionID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to evaluate region")
		}
		return
	}

	RespondJSON(w, http.StatusOK, evaluation)
}

func (a *API) EvaluateRegions(w http.ResponseWriter, r *http.Request) {
	/*ctx := r.Context()
	userID := UserID(ctx)

	evaluations, err := a.evaluationSvc.EvaluateRegions(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondJSON(w, http.StatusNotFound, "region not found")
		default:
			a.logger.Error("failed to evaluate regions", "userId", userID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to evaluate regions")
		}
		return
	}

	RespondJSON(w, http.StatusOK, evaluations)*/
}

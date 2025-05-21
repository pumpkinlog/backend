package api

import (
	"net/http"

	"github.com/pumpkinlog/backend/internal/apiutil"
)

type EvaluateRegionRequest struct {
	RegionID string `json:"regionId" validate:"required,min=2,max=5"`
}

func (a *API) EvaluateRegion(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	params := &EvaluateRegionRequest{
		RegionID: r.PathValue("regionId"),
	}

	if ok := apiutil.Validate(w, params); !ok {
		return
	}

	res, err := a.evaluationSvc.EvaluateRegion(r.Context(), userID, params.RegionID)
	if err != nil {
		a.logger.Error("failed to analyze region", "regionId", params.RegionID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to analyze region")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, res)
}

type EvaluateRegionsRequest struct {
	RegionID []string `json:"regionIds" validate:"dive,min=2,max=5"`
}

func (a *API) EvaluateRegions(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	var params EvaluateRegionsRequest
	if ok := apiutil.ParseJSON(w, r, params); !ok {
		return
	}

	if ok := apiutil.Validate(w, &params); !ok {
		return
	}

	res, err := a.evaluationSvc.EvaluateRegions(r.Context(), userID, params.RegionID)
	if err != nil {
		a.logger.Error("failed to analyze regions", "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to analyze regions")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, res)
}

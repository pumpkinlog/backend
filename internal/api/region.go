package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pumpkinlog/backend/internal/apiutil"
	"github.com/pumpkinlog/backend/internal/domain"
)

type GetRegionRequest struct {
	RegionID string `json:"regionId" validate:"required,min=2,max=5"`
}

func (a *API) GetRegion(w http.ResponseWriter, r *http.Request) {

	params := &GetRegionRequest{
		RegionID: r.PathValue("regionId"),
	}

	if ok := apiutil.Validate(w, params); !ok {
		return
	}

	region, err := a.regionRepo.GetByID(r.Context(), params.RegionID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			apiutil.RespondError(w, http.StatusNotFound, "region not found")
			return
		}

		a.logger.Error("failed to get region", "regionId", params.RegionID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to get region")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, region)
}

type ListRegionsRequest struct {
	RegionIDs []string `json:"regionId" validate:"dive,min=2,max=5"`
	Page      *int     `json:"page" validate:"omitempty,min=1"`
	Limit     *int     `json:"limit" validate:"omitempty,min=1,max=100"`
}

func (a *API) ListRegions(w http.ResponseWriter, r *http.Request) {

	page, err := apiutil.ParseIntPtr(r.URL.Query().Get("page"))
	if err != nil {
		apiutil.RespondError(w, http.StatusBadRequest, fmt.Sprintf("cannot parse page: %s", err.Error()))
		return
	}

	limit, err := apiutil.ParseIntPtr(r.URL.Query().Get("limit"))
	if err != nil {
		apiutil.RespondError(w, http.StatusBadRequest, fmt.Sprintf("cannot parse limit: %s", err.Error()))
		return
	}

	params := &ListRegionsRequest{
		RegionIDs: r.URL.Query()["regionId"],
		Page:      page,
		Limit:     limit,
	}

	if ok := apiutil.Validate(w, params); !ok {
		return
	}

	filter := &domain.RegionFilter{
		RegionIDs: params.RegionIDs,
		Page:      params.Page,
		Limit:     params.Limit,
	}

	regions, err := a.regionRepo.List(r.Context(), filter)
	if err != nil {
		a.logger.Error("failed to list regions", "regionIds", params.RegionIDs, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to list regions")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, regions)
}

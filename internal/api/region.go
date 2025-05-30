package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetRegion(w http.ResponseWriter, r *http.Request) {

	regionID := r.PathValue("regionId")

	region, err := a.regionRepo.GetByID(r.Context(), regionID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondError(w, http.StatusNotFound, "region not found")
		default:
			a.logger.Error("failed to get region", "regionId", regionID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to get region")
		}
		return
	}

	RespondJSON(w, http.StatusOK, region)
}

func (a *API) ListRegions(w http.ResponseWriter, r *http.Request) {

	regionIDs := r.URL.Query()["regionId"]

	page := 1
	limit := PaginationDefaultLimit

	if v := r.URL.Query().Get("page"); v != "" {
		num, err := strconv.Atoi(v)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid page")
			return
		}
		page = num
	}

	if v := r.URL.Query().Get("limit"); v != "" {
		num, err := strconv.Atoi(v)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = num
	}

	if page < 1 {
		page = 1
	}

	if limit > PaginationMaxLimit {
		RespondError(w, http.StatusBadRequest, fmt.Sprintf("limit cannot be greater than: %d", PaginationMaxLimit))
		return
	}

	filter := &domain.RegionFilter{
		RegionIDs: regionIDs,
		Page:      &page,
		Limit:     &limit,
	}

	regions, err := a.regionRepo.List(r.Context(), filter)
	if err != nil {
		a.logger.Error("failed to list regions", "regionIds", regionIDs, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to list regions")
		return
	}

	RespondJSON(w, http.StatusOK, regions)
}

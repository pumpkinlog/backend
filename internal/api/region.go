package api

import (
	"errors"
	"net/http"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetRegion(w http.ResponseWriter, r *http.Request) {
	regionID := domain.RegionID(r.PathValue("regionId"))

	region, err := a.regionSvc.GetByID(r.Context(), regionID)
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
	regionIDs := make([]domain.RegionID, 0)
	for _, rid := range r.URL.Query()["regionId"] {
		regionIDs = append(regionIDs, domain.RegionID(rid))
	}

	filter := &domain.RegionFilter{
		RegionIDs: regionIDs,
	}

	regions, err := a.regionSvc.List(r.Context(), filter)
	if err != nil {
		a.logger.Error("failed to list regions", "regionIds", regionIDs, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to list regions")
		return
	}

	RespondJSON(w, http.StatusOK, regions)
}

package api

import (
	"errors"
	"net/http"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetCondition(w http.ResponseWriter, r *http.Request) {
	conditionID := domain.Code(r.PathValue("conditionId"))

	condition, err := a.conditionSvc.GetByID(r.Context(), conditionID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondError(w, http.StatusNotFound, "condition not found")
		default:
			a.logger.Error("failed to get condition", "conditionId", conditionID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to get condition")
		}
		return
	}

	RespondJSON(w, http.StatusOK, condition)
}

func (a *API) ListConditions(w http.ResponseWriter, r *http.Request) {
	regionIDs := make([]domain.RegionID, 0)
	for _, rid := range r.URL.Query()["regionId"] {
		regionIDs = append(regionIDs, domain.RegionID(rid))
	}

	filter := &domain.ConditionFilter{
		RegionIDs: regionIDs,
	}

	conditions, err := a.conditionSvc.List(r.Context(), filter)
	if err != nil {
		a.logger.Error("failed to list conditions", "regionIds", regionIDs, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to list conditions")
		return
	}

	RespondJSON(w, http.StatusOK, conditions)
}

package api

import (
	"errors"
	"net/http"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetRule(w http.ResponseWriter, r *http.Request) {
	ruleID := domain.Code(r.PathValue("ruleId"))

	rule, err := a.ruleSvc.GetByID(r.Context(), ruleID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondError(w, http.StatusNotFound, "rule not found")
		default:
			a.logger.Error("failed to get rule", "ruleId", ruleID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to get rule")
		}
		return
	}

	RespondJSON(w, http.StatusOK, rule)
}

func (a *API) ListRules(w http.ResponseWriter, r *http.Request) {
	regionIDs := make([]domain.RegionID, 0)
	for _, rid := range r.URL.Query()["regionId"] {
		regionIDs = append(regionIDs, domain.RegionID(rid))
	}

	filter := &domain.RuleFilter{
		RegionIDs: regionIDs,
	}

	rules, err := a.ruleSvc.List(r.Context(), filter)
	if err != nil {
		a.logger.Error("failed to list rules", "regionIds", regionIDs, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to list rules")
		return
	}

	RespondJSON(w, http.StatusOK, rules)
}

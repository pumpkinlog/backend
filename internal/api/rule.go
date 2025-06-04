package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetRule(w http.ResponseWriter, r *http.Request) {

	var err error
	var ruleID int64
	if ruleID, err = strconv.ParseInt(r.PathValue("ruleId"), 10, 64); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid rule ID")
		return
	}

	rule, err := a.ruleRepo.GetByID(r.Context(), ruleID)
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

	regionIDs := r.URL.Query()["regionId"]

	filter := &domain.RuleFilter{
		RegionIDs: regionIDs,
	}

	rules, err := a.ruleRepo.List(r.Context(), filter)
	if err != nil {
		a.logger.Error("failed to list rules", "regionIds", regionIDs, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to list rules")
		return
	}

	RespondJSON(w, http.StatusOK, rules)
}

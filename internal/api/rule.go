package api

import (
	"errors"
	"net/http"

	"github.com/pumpkinlog/backend/internal/apiutil"
	"github.com/pumpkinlog/backend/internal/domain"
)

type GetRuleRequest struct {
	RuleID string `json:"ruleId" validate:"required,uuid"`
}

func (a *API) GetRule(w http.ResponseWriter, r *http.Request) {

	params := &GetRuleRequest{
		RuleID: r.PathValue("ruleId"),
	}

	if ok := apiutil.Validate(w, params); !ok {
		return
	}

	rule, err := a.ruleRepo.GetByID(r.Context(), params.RuleID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			apiutil.RespondError(w, http.StatusNotFound, "rule not found")
			return
		}

		a.logger.Error("failed to get rule", "ruleId", params.RuleID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to get rule")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, rule)
}

type ListRulesRequest struct {
	// @TODO: Perhaps require regionIDs to be set to prevent listing all rules?
	RegionID []string `json:"regionId" validate:"dive,min=2,max=5"`
}

func (a *API) ListRules(w http.ResponseWriter, r *http.Request) {

	params := &ListRulesRequest{
		RegionID: r.URL.Query()["regionId"],
	}

	if ok := apiutil.Validate(w, params); !ok {
		return
	}

	filter := &domain.RuleFilter{
		RegionIDs: params.RegionID,
	}

	rules, err := a.ruleRepo.List(r.Context(), filter)
	if err != nil {
		a.logger.Error("failed to list rules", "regionId", params.RegionID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to list rules")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, rules)
}

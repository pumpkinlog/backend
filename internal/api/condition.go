package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pumpkinlog/backend/internal/apiutil"
	"github.com/pumpkinlog/backend/internal/domain"
)

type GetConditionRequest struct {
	ConditionID string `json:"conditionId" validate:"required,uuid"`
}

func (a *API) GetCondition(w http.ResponseWriter, r *http.Request) {

	params := &GetConditionRequest{
		ConditionID: r.PathValue("conditionId"),
	}

	if ok := apiutil.Validate(w, params); !ok {
		return
	}

	condition, err := a.conditionRepo.GetByID(r.Context(), params.ConditionID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			apiutil.RespondError(w, http.StatusNotFound, "condition not found")
			return
		}

		a.logger.Error("failed to get condition", "conditionId", params.ConditionID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to get condition")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, condition)
}

type ListConditionsRequest struct {
	// @TODO: Perhaps require ruleIDs to be set to prevent listing all conditions
	RuleIDs []string `json:"ruleId" validate:"dive,uuid"`
	Page    *int     `json:"page" validate:"omitempty,min=1"`
	Limit   *int     `json:"limit" validate:"omitempty,min=1,max=100"`
}

func (a *API) ListConditions(w http.ResponseWriter, r *http.Request) {

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

	params := &ListConditionsRequest{
		RuleIDs: r.URL.Query()["ruleId"],
		Page:    page,
		Limit:   limit,
	}

	if ok := apiutil.Validate(w, params); !ok {
		return
	}

	filter := &domain.ConditionFilter{
		RuleIDs: params.RuleIDs,
	}

	conditions, err := a.conditionRepo.List(r.Context(), filter)
	if err != nil {
		a.logger.Error("failed to list conditions", "ruleIds", params.RuleIDs, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to list conditions")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, conditions)
}

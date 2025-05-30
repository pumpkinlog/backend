package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetCondition(w http.ResponseWriter, r *http.Request) {

	conditionID := r.PathValue("conditionId")

	condition, err := a.conditionRepo.GetByID(r.Context(), conditionID)
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

	filter := &domain.ConditionFilter{
		Page:  &page,
		Limit: &limit,
	}

	conditions, err := a.conditionRepo.List(r.Context(), filter)
	if err != nil {
		a.logger.Error("failed to list conditions", "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to list conditions")
		return
	}

	RespondJSON(w, http.StatusOK, conditions)
}

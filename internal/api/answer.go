package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetAnswer(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r)

	var err error
	var conditionID int64
	if conditionID, err = strconv.ParseInt(r.PathValue("conditionId"), 10, 64); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid condition ID")
		return
	}

	answer, err := a.answerRepo.GetByID(r.Context(), userID, conditionID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondJSON(w, http.StatusNotFound, "answer not found")
		default:
			a.logger.Error("failed to get answer", "userId", userID, "conditionId", conditionID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to get answer")
		}
		return
	}

	RespondJSON(w, http.StatusOK, answer)
}

type SubmitAnswerRequest struct {
	ConditionID int64 `json:"conditionId"`
	Value       any   `json:"value"`
}

func (a *API) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r)

	var params SubmitAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondError(w, http.StatusBadRequest, "malformed request body")
		return
	}
	defer r.Body.Close()

	if err := a.answerSvc.CreateOrUpdate(r.Context(), userID, params.ConditionID, params.Value); err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			a.logger.Error("failed to create or update answer", "userId", userID, "conditionId", params.ConditionID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to create or update answer")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r)

	var err error
	var conditionID int64
	if conditionID, err = strconv.ParseInt(r.PathValue("conditionId"), 10, 64); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid condition ID")
		return
	}

	if err := a.answerRepo.Delete(r.Context(), userID, conditionID); err != nil {
		a.logger.Error("failed to delete answer", "userId", userID, "conditionId", conditionID, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to delete answer")
		return
	}

	w.WriteHeader(http.StatusOK)
}

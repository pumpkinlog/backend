package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetAnswer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserID(ctx)
	conditionID := domain.Code(r.PathValue("conditionId"))

	answer, err := a.answerSvc.GetByID(ctx, userID, conditionID)
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
	ConditionID domain.Code `json:"conditionId"`
	Value       any         `json:"value"`
}

func (a *API) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserID(ctx)

	var params SubmitAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondError(w, http.StatusBadRequest, "malformed request body")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if err := a.answerSvc.CreateOrUpdate(ctx, userID, params.ConditionID, params.Value); err != nil {
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
	ctx := r.Context()
	userID := UserID(ctx)
	conditionID := domain.Code(r.PathValue("conditionId"))

	if err := a.answerSvc.Delete(ctx, userID, conditionID); err != nil {
		a.logger.Error("failed to delete answer", "userId", userID, "conditionId", conditionID, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to delete answer")
		return
	}

	w.WriteHeader(http.StatusOK)
}

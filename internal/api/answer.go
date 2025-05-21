package api

import (
	"net/http"

	"github.com/pumpkinlog/backend/internal/apiutil"
	"github.com/pumpkinlog/backend/internal/domain"
)

type SubmitAnswerRequest struct {
	ConditionID string `json:"conditionId"`
	Value       any    `json:"value"`
}

func (a *API) SubmitAnswer(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	var params SubmitAnswerRequest
	if ok := apiutil.ParseJSON(w, r, &params); !ok {
		return
	}

	if ok := apiutil.Validate(w, &params); !ok {
		return
	}

	answer := &domain.Answer{
		UserID:      userID,
		ConditionID: params.ConditionID,
		Value:       params.Value,
	}

	if err := a.answerRepo.CreateOrUpdate(r.Context(), answer); err != nil {
		a.logger.Error("failed to create or update answer", "userId", userID, "conditionId", params.ConditionID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to create or update answer")
		return
	}

	w.WriteHeader(http.StatusOK)
}

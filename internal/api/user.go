package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/pumpkinlog/backend/internal/apiutil"
	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetUser(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	user, err := a.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			apiutil.RespondError(w, http.StatusNotFound, "user not found")
			return
		}

		a.logger.Error("failed to get user", "userId", userID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, user)
}

type CreateUserRequest struct {
	UserID string `json:"userId" validate:"required,uuid"`
}

func (a *API) CreateUser(w http.ResponseWriter, r *http.Request) {

	var params CreateUserRequest
	if ok := apiutil.ParseJSON(w, r, &params); !ok {
		return
	}

	if ok := apiutil.Validate(w, &params); !ok {
		return
	}

	user := &domain.User{
		ID:              params.UserID,
		FavoriteRegions: make([]string, 0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := a.userRepo.Create(r.Context(), user); err != nil {
		a.logger.Error("failed to create user", "userId", params.UserID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	w.WriteHeader(http.StatusOK)
}

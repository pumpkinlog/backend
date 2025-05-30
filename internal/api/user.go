package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r)

	user, err := a.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondError(w, http.StatusNotFound, "user not found")
		default:
			a.logger.Error("failed to get user", "userId", userID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to get user")
		}
		return
	}

	RespondJSON(w, http.StatusOK, user)
}

type CreateUserRequest struct {
	UserID string `json:"userId"`
}

func (a *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	// @TODO: validate webhook signature

	var params CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondError(w, http.StatusBadRequest, "malformed request body")
		return
	}
	defer r.Body.Close()

	if params.UserID == "" {
		RespondError(w, http.StatusBadRequest, "userId is required")
		return
	}

	if err := a.userSvc.Create(r.Context(), params.UserID); err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			a.logger.Error("failed to create user", "userId", params.UserID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to create user")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type UpdateUserRequest struct {
	FavoriteRegions []string `json:"favoriteRegions"`
	WantResidency   []string `json:"wantResidency"`
}

func (a *API) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r)

	var params UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondError(w, http.StatusBadRequest, "malformed request body")
		return
	}
	defer r.Body.Close()

	if err := a.userSvc.Update(r.Context(), userID, params.FavoriteRegions, params.WantResidency); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondError(w, http.StatusNotFound, "user not found")
		case errors.Is(err, domain.ErrValidation):
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			a.logger.Error("failed to update user", "userId", userID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to update user")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserID(ctx)

	user, err := a.userSvc.GetByID(ctx, userID)
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
	FavoriteRegions []domain.RegionID `json:"favoriteRegions"`
	WantResidency   []domain.RegionID `json:"wantResidency"`
}

func (a *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	// @TODO: validate webhook signature

	var params CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondError(w, http.StatusBadRequest, "malformed request body")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if err := a.userSvc.Create(r.Context(), params.FavoriteRegions, params.WantResidency); err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			a.logger.Error("failed to create user", "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to create user")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type UpdateUserRequest struct {
	FavoriteRegions []domain.RegionID `json:"favoriteRegions"`
	WantResidency   []domain.RegionID `json:"wantResidency"`
}

func (a *API) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserID(ctx)

	var params UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondError(w, http.StatusBadRequest, "malformed request body")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if err := a.userSvc.Update(ctx, userID, params.FavoriteRegions, params.WantResidency); err != nil {
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

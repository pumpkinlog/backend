package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pumpkinlog/backend/internal/apiutil"
	"github.com/pumpkinlog/backend/internal/domain"
)

type GetPresenceRequest struct {
	RegionID string    `json:"regionId" validate:"required,min=2,max=5"`
	Date     time.Time `json:"date" validate:"required"`
}

func (a *API) GetPresence(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	date, ok := apiutil.ParseDate(w, r.PathValue("date"))
	if !ok {
		return
	}

	params := &GetPresenceRequest{
		RegionID: r.PathValue("regionId"),
		Date:     date,
	}

	if ok := apiutil.Validate(w, params); !ok {
		return
	}

	presence, err := a.presenceRepo.GetByID(r.Context(), userID, params.RegionID, params.Date)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			apiutil.RespondError(w, http.StatusNotFound, "presence not found")
			return
		}

		apiutil.RespondError(w, http.StatusInternalServerError, "failed to get presence")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, presence)
}

type ListPresenceRequest struct {
	RegionIDs []string   `json:"regionIds" validate:"dive,min=2,max=5"`
	Page      *int       `json:"page" validate:"omitempty,min=1"`
	Limit     *int       `json:"limit" validate:"omitempty,min=1,max=100"`
	Start     *time.Time `json:"start"`
	End       *time.Time `json:"end"`
}

func (a *API) ListPresences(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	start, err := apiutil.ParseDatePtr(r.URL.Query().Get("start"))
	if err != nil {
		apiutil.RespondError(w, http.StatusBadRequest, fmt.Sprintf("cannot parse start: %s", err.Error()))
		return
	}

	end, err := apiutil.ParseDatePtr(r.URL.Query().Get("end"))
	if err != nil {
		apiutil.RespondError(w, http.StatusBadRequest, fmt.Sprintf("cannot parse end: %s", err.Error()))
		return
	}

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

	params := &ListPresenceRequest{
		RegionIDs: r.URL.Query()["regionIds"],
		Start:     start,
		End:       end,
		Page:      page,
		Limit:     limit,
	}

	if ok := apiutil.Validate(w, params); !ok {
		return
	}

	filter := &domain.PresenceFilter{
		RegionIDs: params.RegionIDs,
		Start:     params.Start,
		End:       params.End,
		Page:      params.Page,
		Limit:     params.Limit,
	}

	precences, err := a.presenceRepo.List(r.Context(), userID, filter)
	if err != nil {
		a.logger.Error("failed to list presences", "error", err.Error())
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to list presences")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, precences)
}

type CreatePresencesRequest struct {
	RegionID string    `json:"regionId"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	DeviceID *string   `json:"deviceId"`
}

func (a *API) CreatePresence(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	var params CreatePresencesRequest
	if ok := apiutil.ParseJSON(w, r, &params); !ok {
		return
	}

	if err := a.presenceSvc.Create(r.Context(), userID, params.RegionID, params.DeviceID, params.Start, params.End); err != nil {
		a.logger.Error("failed to create presence", "error", err.Error())
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to create presence")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type DeletePresenceRequest struct {
	RegionID string    `json:"regionId" validate:"required,min=2,max=5"`
	Start    time.Time `json:"start" validate:"required"`
	End      time.Time `json:"end" validate:"required"`
	DeviceID *string   `json:"deviceId" validate:"uuid"`
}

func (a *API) DeletePresence(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	regionID := r.PathValue("regionId")
	if regionID == "" {
		return
	}

	start, ok := apiutil.ParseDate(w, r.PathValue("start"))
	if !ok {
		return
	}

	end, ok := apiutil.ParseDate(w, r.PathValue("end"))
	if !ok {
		return
	}

	if err := a.presenceSvc.Delete(r.Context(), userID, regionID, start, end); err != nil {
		a.logger.Error("failed to delete presence", "error", err.Error())
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to delete presence")
		return
	}

	w.WriteHeader(http.StatusOK)
}

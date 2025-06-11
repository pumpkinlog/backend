package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetPresence(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserID(ctx)
	regionID := domain.RegionID(r.PathValue("regionId"))

	date, err := time.Parse(time.DateOnly, r.PathValue("date"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid date")
		return
	}

	presence, err := a.presenceSvc.GetByID(ctx, userID, regionID, date)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondError(w, http.StatusNotFound, "presence not found")
		default:
			a.logger.Error("failed to get presence", "userId", userID, "regionId", regionID, "date", date, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to get presence")
		}
		return
	}

	RespondJSON(w, http.StatusOK, presence)
}

func (a *API) ListPresences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserID(ctx)

	regionIDs := make([]domain.RegionID, 0)
	for _, rid := range r.URL.Query()["regionId"] {
		regionIDs = append(regionIDs, domain.RegionID(rid))
	}

	var start *time.Time
	var end *time.Time

	if v := r.URL.Query().Get("start"); v != "" {
		t, err := time.Parse(time.DateOnly, v)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid start time")
			return
		}
		start = &t
	}

	if v := r.URL.Query().Get("end"); v != "" {
		t, err := time.Parse(time.DateOnly, v)
		if err != nil {
			RespondError(w, http.StatusBadRequest, "invalid end time")
			return
		}
		end = &t
	}

	filter := &domain.PresenceFilter{
		RegionIDs: regionIDs,
		Start:     start,
		End:       end,
	}

	precences, err := a.presenceSvc.List(ctx, userID, filter)
	if err != nil {
		a.logger.Error("failed to list presences", "userId", userID, "regionIds", regionIDs, "start", start, "end", end, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to list presences")
		return
	}

	RespondJSON(w, http.StatusOK, precences)
}

type CreatePresencesRequest struct {
	RegionID domain.RegionID `json:"regionId"`
	Start    string          `json:"start"`
	End      string          `json:"end"`
	DeviceID *int64          `json:"deviceId"`
}

func (a *API) CreatePresence(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserID(ctx)

	var params CreatePresencesRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondError(w, http.StatusBadRequest, "malformed request body")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	start, err := time.Parse(time.DateOnly, params.Start)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid start time")
		return
	}

	end, err := time.Parse(time.DateOnly, params.End)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid end time")
		return
	}

	if err := a.presenceSvc.Create(ctx, userID, params.RegionID, params.DeviceID, start, end); err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
		default:
			a.logger.Error("failed to create presence", "userId", userID, "regionId", params.RegionID, "deviceId", params.DeviceID, "start", start, "end", end, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to create presence")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) DeletePresence(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := UserID(ctx)
	regionID := domain.RegionID(r.URL.Query().Get("regionId"))

	start, err := time.Parse(time.DateOnly, r.URL.Query().Get("start"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid start time")
		return
	}

	end, err := time.Parse(time.DateOnly, r.URL.Query().Get("end"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "invalid end time")
		return
	}

	if err := a.presenceSvc.Delete(ctx, userID, regionID, start, end); err != nil {
		a.logger.Error("failed to delete presence", "userId", userID, "regionId", regionID, "start", start, "end", end, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to delete presence")
		return
	}

	w.WriteHeader(http.StatusOK)
}

package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/pumpkinlog/backend/internal/domain"
)

func (a *API) GetDevice(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r)

	var err error
	var deviceID int64
	if deviceID, err = strconv.ParseInt(r.PathValue("deviceId"), 10, 64); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid device ID")
		return
	}

	device, err := a.deviceSvc.GetByID(r.Context(), userID, deviceID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondError(w, http.StatusNotFound, "device not found")
		default:
			a.logger.Error("failed to get device", "deviceId", deviceID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to get device")
		}
		return
	}

	RespondJSON(w, http.StatusOK, device)
}

func (a *API) ListDevices(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r)

	devices, err := a.deviceSvc.List(r.Context(), userID)
	if err != nil {
		a.logger.Error("failed to list devices", "userId", userID, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to list devices")
		return
	}

	RespondJSON(w, http.StatusOK, devices)
}

type CreateDeviceRequest struct {
	Name     string  `json:"name"`
	Platform string  `json:"platform"`
	Model    string  `json:"model"`
	Token    *string `json:"token"`
	Active   bool    `json:"active"`
}

func (a *API) CreateDevice(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r)

	var params CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondError(w, http.StatusBadRequest, "malformed request body")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if err := a.deviceSvc.Create(r.Context(), userID, params.Name, params.Platform, params.Model); err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			a.logger.Error("failed to create device", "userId", userID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to create device")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type UpdateDeviceRequest struct {
	DeviceID int64  `json:"deviceId"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	Active   bool   `json:"active"`
}

func (a *API) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r)

	var params UpdateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondError(w, http.StatusBadRequest, "malformed request body")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if err := a.deviceSvc.Update(r.Context(), userID, params.DeviceID, params.Name, params.Token, params.Active); err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			RespondError(w, http.StatusNotFound, "device not found")
		case errors.Is(err, domain.ErrValidation):
			RespondError(w, http.StatusBadRequest, err.Error())
		default:
			a.logger.Error("failed to update device", "userId", userID, "error", err)
			RespondError(w, http.StatusInternalServerError, "failed to updated device")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	userID := UserID(r)

	var err error
	var deviceID int64
	if deviceID, err = strconv.ParseInt(r.PathValue("deviceId"), 10, 64); err != nil {
		RespondError(w, http.StatusBadRequest, "invalid device ID")
		return
	}

	if err := a.deviceSvc.Delete(r.Context(), userID, deviceID); err != nil {
		a.logger.Error("failed to delete device", "deviceId", deviceID, "error", err)
		RespondError(w, http.StatusInternalServerError, "failed to delete device")
		return
	}

	w.WriteHeader(http.StatusOK)
}

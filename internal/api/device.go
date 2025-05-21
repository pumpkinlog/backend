package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/pumpkinlog/backend/internal/apiutil"
	"github.com/pumpkinlog/backend/internal/domain"
)

type GetDeviceRequest struct {
	DeviceID string `json:"deviceId" validate:"required,uuid"`
}

func (a *API) GetDevice(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	params := &GetDeviceRequest{
		DeviceID: r.PathValue("deviceId"),
	}

	device, err := a.deviceRepo.GetByID(r.Context(), userID, params.DeviceID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			apiutil.RespondError(w, http.StatusNotFound, "device not found")
			return
		}

		a.logger.Error("failed to get device", "deviceId", params.DeviceID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to get device")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, device)
}

func (a *API) ListDevices(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	devices, err := a.deviceRepo.List(r.Context(), userID)
	if err != nil {
		a.logger.Error("failed to list devices", "userId", userID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to list devices")
		return
	}

	apiutil.RespondJSON(w, http.StatusOK, devices)
}

type CreateDeviceRequest struct {
	Name     string  `json:"name"`
	Platform string  `json:"platform"`
	Model    string  `json:"model"`
	Token    *string `json:"token"`
	Active   bool    `json:"active"`
}

type CreateDeviceResponse struct {
	DeviceID string `json:"deviceId"`
}

func (a *API) CreateDevice(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	var params CreateDeviceRequest
	if ok := apiutil.ParseJSON(w, r, &params); !ok {
		return
	}

	if ok := apiutil.Validate(w, &params); !ok {
		return
	}

	device := &domain.Device{
		UserID:    userID,
		Name:      params.Name,
		Platform:  domain.Platform(params.Platform),
		Model:     params.Model,
		Token:     params.Token,
		Active:    params.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.deviceRepo.Create(r.Context(), device); err != nil {
		a.logger.Error("failed to create device", "userId", userID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to create device")
		return
	}

	res := CreateDeviceResponse{
		DeviceID: device.ID,
	}

	apiutil.RespondJSON(w, http.StatusCreated, res)
}

type UpdateDeviceRequest struct {
	DeviceID string  `json:"deviceId"`
	Name     *string `json:"name"`
	Token    *string `json:"token"`
	Active   *bool   `json:"active"`
}

func (a *API) UpdateDevice(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	var params UpdateDeviceRequest
	if ok := apiutil.ParseJSON(w, r, &params); !ok {
		return
	}

	if ok := apiutil.Validate(w, &params); !ok {
		return
	}

	device, err := a.deviceRepo.GetByID(r.Context(), userID, params.DeviceID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			apiutil.RespondError(w, http.StatusNotFound, "device not found")
			return
		}

		a.logger.Error("failed to get device", "deviceId", params.DeviceID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to get device")
		return
	}

	if params.Name != nil {
		device.Name = *params.Name
	}

	if params.Token != nil {
		device.Token = params.Token
	}

	if params.Active != nil {
		device.Active = *params.Active
	}

	if err := a.deviceRepo.Update(r.Context(), device); err != nil {
		a.logger.Error("failed to update device", "deviceId", params.DeviceID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to update device")
		return
	}

	w.WriteHeader(http.StatusOK)
}

type DeleteDeviceRequest struct {
	DeviceID string `json:"deviceId" validate:"required,uuid"`
}

func (a *API) DeleteDevice(w http.ResponseWriter, r *http.Request) {

	userID, ok := apiutil.UserID(w, r)
	if !ok {
		return
	}

	params := &DeleteDeviceRequest{
		DeviceID: r.PathValue("deviceId"),
	}

	if ok := apiutil.Validate(w, params); !ok {
		return
	}

	if err := a.deviceRepo.Delete(r.Context(), userID, params.DeviceID); err != nil {
		a.logger.Error("failed to delete device", "deviceId", params.DeviceID, "error", err)
		apiutil.RespondError(w, http.StatusInternalServerError, "failed to delete device")
		return
	}

	w.WriteHeader(http.StatusOK)
}

package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateDevice(t *testing.T) {
	timestamp := time.Now()
	token := "sometoken"
	device := Device{
		ID:        1,
		UserID:    1,
		Name:      "My Device",
		Platform:  PlatformIOS,
		Model:     "iPhone 15",
		Token:     &token,
		Active:    true,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	tests := []struct {
		name    string
		modify  func(d Device) Device
		wantErr error
	}{
		{
			name:   "valid device",
			modify: func(d Device) Device { return d },
		},
		{
			name: "invalid user ID",
			modify: func(d Device) Device {
				d.UserID = 0
				return d
			},
			wantErr: ValidationError("user ID is required"),
		},
		{
			name: "missing name",
			modify: func(d Device) Device {
				d.Name = ""
				return d
			},
			wantErr: ValidationError("name is required"),
		},
		{
			name: "invalid platform",
			modify: func(d Device) Device {
				d.Platform = "windows"
				return d
			},
			wantErr: ValidationError("platform is invalid"),
		},
		{
			name: "missing model",
			modify: func(d Device) Device {
				d.Model = ""
				return d
			},
			wantErr: ValidationError("model is required"),
		},
		{
			name: "missing created at",
			modify: func(d Device) Device {
				d.CreatedAt = time.Time{}
				return d
			},
			wantErr: ValidationError("created at timestamp is invalid"),
		},
		{
			name: "missing updated at",
			modify: func(d Device) Device {
				d.UpdatedAt = time.Time{}
				return d
			},
			wantErr: ValidationError("updated at timestamp is invalid"),
		},
		{
			name: "created at in the future",
			modify: func(d Device) Device {
				d.CreatedAt = time.Now().Add(time.Minute)
				return d
			},
			wantErr: ValidationError("created at timestamp cannot be in the future"),
		},
		{
			name: "updated at in the future",
			modify: func(d Device) Device {
				d.UpdatedAt = time.Now().Add(time.Minute)
				return d
			},
			wantErr: ValidationError("updated at timestamp cannot be in the future"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dev := tc.modify(device)
			err := dev.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

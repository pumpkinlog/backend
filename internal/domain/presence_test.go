package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidatePresence(t *testing.T) {
	timestamp := time.Now()
	deviceID := "device-123"
	presence := Presence{
		UserID:    1,
		RegionID:  "JE",
		Date:      timestamp,
		DeviceID:  &deviceID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	tests := []struct {
		name    string
		modify  func(p Presence) Presence
		wantErr error
	}{
		{
			name:   "valid presence",
			modify: func(p Presence) Presence { return p },
		},
		{
			name: "invalid user ID",
			modify: func(p Presence) Presence {
				p.UserID = -1
				return p
			},
			wantErr: ValidationError("user ID is required"),
		},
		{
			name: "invalid region ID",
			modify: func(p Presence) Presence {
				p.RegionID = ""
				return p
			},
			wantErr: ValidationError("region ID is required"),
		},
		{
			name: "date in the future",
			modify: func(p Presence) Presence {
				p.Date = time.Now().Add(time.Minute)
				return p
			},
			wantErr: ValidationError("date cannot be in the future"),
		},
		{
			name: "device ID empty string",
			modify: func(p Presence) Presence {
				empty := ""
				p.DeviceID = &empty
				return p
			},
			wantErr: ValidationError("device ID cannot be empty"),
		},
		{
			name: "nil device ID is allowed",
			modify: func(p Presence) Presence {
				p.DeviceID = nil
				return p
			},
		},
		{
			name: "missing created at",
			modify: func(p Presence) Presence {
				p.CreatedAt = time.Time{}
				return p
			},
			wantErr: ValidationError("created at timestamp is required"),
		},
		{
			name: "missing updated at",
			modify: func(p Presence) Presence {
				p.UpdatedAt = time.Time{}
				return p
			},
			wantErr: ValidationError("updated at timestamp is required"),
		},
		{
			name: "created at in the future",
			modify: func(p Presence) Presence {
				p.CreatedAt = time.Now().Add(time.Minute)
				return p
			},
			wantErr: ValidationError("created at timestamp cannot be in the future"),
		},
		{
			name: "updated at in the future",
			modify: func(p Presence) Presence {
				p.UpdatedAt = time.Now().Add(time.Minute)
				return p
			},
			wantErr: ValidationError("updated at timestamp cannot be in the future"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.modify(presence)
			err := p.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

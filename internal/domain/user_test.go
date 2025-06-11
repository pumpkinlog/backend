package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateUser(t *testing.T) {
	timestamp := time.Now()
	user := User{
		ID:              1,
		FavoriteRegions: []RegionID{"JE", "GB"},
		WantResidency:   []RegionID{"US", "CA"},
		CreatedAt:       timestamp,
		UpdatedAt:       timestamp,
	}

	tests := []struct {
		name    string
		modify  func(u User) User
		wantErr error
	}{
		{
			name:   "valid user",
			modify: func(u User) User { return u },
		},
		{
			name: "invalid user ID",
			modify: func(u User) User {
				u.ID = -1
				return u
			},
			wantErr: ValidationError("user ID is required"),
		},
		{
			name: "too many favorite regions",
			modify: func(u User) User {
				u.FavoriteRegions = make([]RegionID, maxRegions+1)
				for i := range u.FavoriteRegions {
					u.FavoriteRegions[i] = "JE"
				}
				return u
			},
			wantErr: ValidationError("favorite regions cannot be greater than %d", maxRegions),
		},
		{
			name: "too many want residency regions",
			modify: func(u User) User {
				u.WantResidency = make([]RegionID, maxRegions+1)
				for i := range u.WantResidency {
					u.WantResidency[i] = "US"
				}
				return u
			},
			wantErr: ValidationError(" want residency cannot be greater than %d", maxRegions),
		},
		{
			name: "missing created at",
			modify: func(u User) User {
				u.CreatedAt = time.Time{}
				return u
			},
			wantErr: ValidationError("created at timestamp is required"),
		},
		{
			name: "missing updated at",
			modify: func(u User) User {
				u.UpdatedAt = time.Time{}
				return u
			},
			wantErr: ValidationError("updated at timestamp is required"),
		},
		{
			name: "created at in the future",
			modify: func(u User) User {
				u.CreatedAt = time.Now().Add(time.Minute)
				return u
			},
			wantErr: ValidationError("created at timestamp cannot be in the future"),
		},
		{
			name: "updated at in the future",
			modify: func(u User) User {
				u.UpdatedAt = time.Now().Add(time.Minute)
				return u
			},
			wantErr: ValidationError("updated at timestamp cannot be in the future"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user := tc.modify(user)
			err := user.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

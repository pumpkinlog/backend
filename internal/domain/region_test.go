package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRegion(t *testing.T) {
	baseID := RegionID("JE")
	parentID := RegionID("GB")

	region := Region{
		ID:             baseID,
		ParentRegionID: &parentID,
		Name:           "Jersey",
		Type:           RegionTypeCountry,
		Continent:      ContinentEurope,
		YearStartMonth: 1,
		YearStartDay:   1,
		LatLng:         [2]float64{49.21, -2.13},
		Sources:        []Source{{Name: "Wikipedia", URL: "https://wikipedia.org"}},
	}

	tests := []struct {
		name    string
		modify  func(r Region) Region
		wantErr error
	}{
		{
			name:   "valid region",
			modify: func(r Region) Region { return r },
		},
		{
			name: "invalid ID",
			modify: func(r Region) Region {
				r.ID = ""
				return r
			},
			wantErr: ValidationError("region ID is required"),
		},
		{
			name: "invalid ParentRegionID",
			modify: func(r Region) Region {
				invalidParent := RegionID("")
				r.ParentRegionID = &invalidParent
				return r
			},
			wantErr: ValidationError("region ID is required"),
		},
		{
			name: "missing name",
			modify: func(r Region) Region {
				r.Name = ""
				return r
			},
			wantErr: ValidationError("name is required"),
		},
		{
			name: "invalid type",
			modify: func(r Region) Region {
				r.Type = "invalid"
				return r
			},
			wantErr: ValidationError("type is required"),
		},
		{
			name: "invalid continent",
			modify: func(r Region) Region {
				r.Continent = "InvalidContinent"
				return r
			},
			wantErr: ValidationError("continent is required"),
		},
		{
			name: "year start month too low",
			modify: func(r Region) Region {
				r.YearStartMonth = 0
				return r
			},
			wantErr: ValidationError("year start month must be between 1-12"),
		},
		{
			name: "year start month too high",
			modify: func(r Region) Region {
				r.YearStartMonth = 13
				return r
			},
			wantErr: ValidationError("year start month must be between 1-12"),
		},
		{
			name: "year start day too low",
			modify: func(r Region) Region {
				r.YearStartDay = 0
				return r
			},
			wantErr: ValidationError("year start day must be between 1-31"),
		},
		{
			name: "year start day too high",
			modify: func(r Region) Region {
				r.YearStartDay = 32
				return r
			},
			wantErr: ValidationError("year start day must be between 1-31"),
		},
		{
			name: "empty sources",
			modify: func(r Region) Region {
				r.Sources = nil
				return r
			},
			wantErr: ValidationError("sources cannot be empty"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.modify(region)
			err := r.Validate()
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

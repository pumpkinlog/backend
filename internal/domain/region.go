package domain

import (
	"context"
	"time"
)

type (
	RegionType string
	Continent  string
)

const (
	RegionTypeCountry  RegionType = "country"
	RegionTypeProvince RegionType = "province"
	RegionTypeZone     RegionType = "zone"

	ContinentAfrica       Continent = "Africa"
	ContinentAntarctica   Continent = "Antarctica"
	ContinentAsia         Continent = "Asia"
	ContinentEurope       Continent = "Europe"
	ContinentNorthAmerica Continent = "North America"
	ContinentSouthAmerica Continent = "South America"
	ContinentOceania      Continent = "Oceania"
)

type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Region struct {
	ID             RegionID   `json:"id"`
	ParentRegionID *RegionID  `json:"parentRegionId,omitempty"`
	Name           string     `json:"name"`
	Type           RegionType `json:"type"`
	Continent      Continent  `json:"continent"`
	YearStartMonth time.Month `json:"yearStartMonth"`
	YearStartDay   int        `json:"yearStartDay"`
	LatLng         [2]float64 `json:"latLng"`
	Sources        []Source   `json:"sources"`
}

func (rt RegionType) Valid() bool {
	switch rt {
	case RegionTypeCountry, RegionTypeProvince, RegionTypeZone:
		return true
	default:
		return false
	}
}

func (c Continent) Valid() bool {
	switch c {
	case ContinentAfrica, ContinentAntarctica, ContinentAsia, ContinentEurope, ContinentNorthAmerica, ContinentSouthAmerica, ContinentOceania:
		return true
	default:
		return false
	}
}

func (r *Region) Validate() error {
	if err := r.ID.Validate(); err != nil {
		return err
	}

	if r.ParentRegionID != nil {
		if err := r.ParentRegionID.Validate(); err != nil {
			return err
		}
	}

	if r.Name == "" {
		return ValidationError("name is required")
	}

	if !r.Type.Valid() {
		return ValidationError("type is required")
	}

	if !r.Continent.Valid() {
		return ValidationError("continent is required")
	}

	if r.YearStartMonth < 1 || r.YearStartMonth > 12 {
		return ValidationError("year start month must be between 1-12")
	}

	if r.YearStartDay < 1 || r.YearStartDay > 31 {
		return ValidationError("year start day must be between 1-31")
	}

	if r.Sources == nil {
		return ValidationError("sources cannot be empty")
	}

	return nil
}

type RegionService interface {
	GetByID(ctx context.Context, regionID RegionID) (*Region, error)
	List(ctx context.Context, filter *RegionFilter) ([]*Region, error)
	CreateOrUpdate(ctx context.Context, region *Region) error
}

type RegionFilter struct {
	RegionIDs []RegionID
}

type RegionRepository interface {
	GetByID(ctx context.Context, regionID RegionID) (*Region, error)
	List(ctx context.Context, filter *RegionFilter) ([]*Region, error)
	CreateOrUpdate(ctx context.Context, region *Region) error
}

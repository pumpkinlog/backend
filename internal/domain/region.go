package domain

import (
	"context"
	"fmt"
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

type Region struct {
	ID             string     `json:"id"`
	ParentRegionID *string    `json:"parentRegionId,omitempty"`
	Name           string     `json:"name"`
	Type           RegionType `json:"type"`
	Continent      Continent  `json:"continent"`
	YearStartMonth time.Month `json:"yearStartMonth"`
	YearStartDay   int        `json:"yearStartDay"`
	LatLng         [2]float64 `json:"latLng"`
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

	length := len(r.ID)
	if length < 3 || length > 5 {
		return fmt.Errorf("%w: region ID must be between 3-5 characters", ErrValidation)
	}

	if r.ParentRegionID != nil {
		length = len(*r.ParentRegionID)
		if length < 3 || length > 5 {
			return fmt.Errorf("%w: region ID must be between 3-5 characters", ErrValidation)
		}
	}

	if r.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}

	if r.Type.Valid() {
		return fmt.Errorf("%w: type is required", ErrValidation)
	}

	if r.Continent.Valid() {
		return fmt.Errorf("%w: continent is required", ErrValidation)
	}

	if r.YearStartMonth < 1 || r.YearStartMonth > 12 {
		return fmt.Errorf("%w: year start month is invalid", ErrValidation)
	}

	if r.YearStartDay < 1 || r.YearStartDay > 31 {
		return fmt.Errorf("%w: year start day is invalid", ErrValidation)
	}

	return nil
}

type RegionService interface {
	CreateOrUpdate(ctx context.Context, region *Region) error
}

type RegionFilter struct {
	RegionIDs []string
	Page      *int
	Limit     *int
}

type RegionRepository interface {
	GetByID(ctx context.Context, regionID string) (*Region, error)
	List(ctx context.Context, filter *RegionFilter) ([]*Region, error)
	CreateOrUpdate(ctx context.Context, region *Region) error
}

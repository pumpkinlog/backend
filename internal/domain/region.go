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

type RegionFilter struct {
	RegionIDs []string
	Page      *int
	Limit     *int
}

type RegionRepository interface {
	GetByID(ctx context.Context, id string) (*Region, error)
	List(ctx context.Context, filter *RegionFilter) ([]*Region, error)

	CreateOrUpdate(ctx context.Context, region *Region) error
}

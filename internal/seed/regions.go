package seed

import (
	"time"

	"github.com/pumpkinlog/backend/internal/domain"
)

var regions = []domain.Region{
	{
		ID:        "AD",
		Name:      "Andorra",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentEurope,
		LatLng:    [2]float64{42.5063, 1.5211},
	},
	{
		ID:        "AE",
		Name:      "United Arab Emirates",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentAsia,
		LatLng:    [2]float64{23.4241, 53.8478},
	},
	{
		ID:        "AF",
		Name:      "Afghanistan",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentAsia,
		LatLng:    [2]float64{33.9391, 67.7099},
	},
	{
		ID:        "AG",
		Name:      "Antigua and Barbuda",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentNorthAmerica,
		LatLng:    [2]float64{17.0608, -61.7964},
	},
	{
		ID:        "AI",
		Name:      "Anguilla",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentNorthAmerica,
		LatLng:    [2]float64{18.2206, -63.0686},
	},
	{
		ID:        "AL",
		Name:      "Albania",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentEurope,
		LatLng:    [2]float64{41.1533, 20.1683},
	},
	{
		ID:        "AM",
		Name:      "Armenia",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentAsia,
		LatLng:    [2]float64{40.0691, 45.0382},
	},
	{
		ID:        "AO",
		Name:      "Angola",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentAfrica,
		LatLng:    [2]float64{-11.2027, 17.8739},
	},
	{
		ID:        "AR",
		Name:      "Argentina",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentSouthAmerica,
		LatLng:    [2]float64{-38.4161, -63.6167},
	},
	{
		ID:        "AS",
		Name:      "American Samoa",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentOceania,
		LatLng:    [2]float64{-14.2700, -170.1322},
	},
	{
		ID:        "AT",
		Name:      "Austria",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentEurope,
		LatLng:    [2]float64{47.5162, 14.5501},
	},
	{
		ID:        "AU",
		Name:      "Australia",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentOceania,
		LatLng:    [2]float64{-25.2744, 133.7751},
	},
	{
		ID:        "AW",
		Name:      "Aruba",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentNorthAmerica,
		LatLng:    [2]float64{12.5211, -69.9687},
	},
	{
		ID:        "AX",
		Name:      "Åland Islands",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentEurope,
		LatLng:    [2]float64{60.1785, 19.9156},
	},
	{
		ID:        "AZ",
		Name:      "Azerbaijan",
		Type:      domain.RegionTypeCountry,
		Continent: domain.ContinentAsia,
		LatLng:    [2]float64{40.1431, 47.5769},
	},
	{
		ID:             "JE",
		Name:           "Jersey",
		Type:           domain.RegionTypeCountry,
		Continent:      domain.ContinentEurope,
		YearStartMonth: time.January,
		YearStartDay:   1,
		LatLng:         [2]float64{49.2144, -2.1312},
	},
	{
		ID:             "GG",
		Name:           "Guernsey",
		Type:           domain.RegionTypeCountry,
		Continent:      domain.ContinentEurope,
		YearStartMonth: time.January,
		YearStartDay:   1,
		LatLng:         [2]float64{49.4657, -2.5859},
	},
}

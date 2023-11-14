package types

import "strings"

type DataType uint32

const (
	Precipitation = 1 << iota
	AirTemperature
	ExtremeTemperature
	ExtremeWind
	Solar
	Wind
	MorePrecipitation
	MorePhenomena
	SoilTemperature
	WaterEquivalent
	WeatherPhenomena
	CloudType
	Cloudiness
	DewPoint
	Moisture
	Pressure
	Sun
	Visibility
	WindSynopsis
	ClimateIndices
	StationObservations
	StandardFormat
	WindTest
	Soil
	MoreWeatherPhenomena
)

func (dt DataType) String() string {
	switch {
	case dt&Precipitation != 0:
		return "precipitation"
	case dt&AirTemperature != 0:
		return "air_temperature"
	case dt&ExtremeTemperature != 0:
		return "extreme_temperature"
	case dt&ExtremeWind != 0:
		return "extreme_wind"
	case dt&Solar != 0:
		return "solar"
	case dt&Wind != 0:
		return "wind"
	case dt&MorePrecipitation != 0:
		return "more_precip"
	case dt&MorePhenomena != 0:
		return "more_phenomena"
	case dt&SoilTemperature != 0:
		return "soil_temperature"
	case dt&WaterEquivalent != 0:
		return "water_equiv"
	case dt&WeatherPhenomena != 0:
		return "weather_phenomena"
	case dt&CloudType != 0:
		return "cloud_type"
	case dt&Cloudiness != 0:
		return "cloudiness"
	case dt&DewPoint != 0:
		return "dew_point"
	case dt&Moisture != 0:
		return "moisture"
	case dt&Pressure != 0:
		return "pressure"
	case dt&Sun != 0:
		return "sun"
	case dt&Visibility != 0:
		return "visibility"
	case dt&WindSynopsis != 0:
		return "wind_synop"
	case dt&ClimateIndices != 0:
		return "climate_indices"
	case dt&StationObservations != 0:
		return "kl"
	case dt&StandardFormat != 0:
		return "standard_format"
	case dt&WindTest != 0:
		return "wind_test"
	case dt&Soil != 0:
		return "soil"
	case dt&MoreWeatherPhenomena != 0:
		return "more_weather_phenomena"
	}
	return ""
}

// ParseString takes an input string and tries to convert the string into a
// data type.
// If multiple values are supposed to be in a string, the function expects
// the values to be separated by a comma.
// If another separator shall be used, please use the ParseStringWithSeparator
// function
func (dt *DataType) ParseString(s string) {
	dt.ParseStringWithSeperator(s, ",")
}

// ParseStringWithSeperator takes an input string and a separator and
// tries to convert the string into a data type.
// If multiple values are supposed to be in a string, the function expects
// the values to be separated by a comma.
// If a comma is used as separator, please use the ParseString
// function as it expects this by default
func (dt *DataType) ParseStringWithSeperator(s, sep string) {
	var parsedDataTypes DataType
	parts := strings.Split(s, sep)
	for _, part := range parts {
		switch part {
		case "precipitation":
			parsedDataTypes = parsedDataTypes | Precipitation
		case "air_temperature":
			parsedDataTypes = parsedDataTypes | AirTemperature
		case "extreme_temperature":
			parsedDataTypes = parsedDataTypes | ExtremeTemperature
		case "extreme_wind":
			parsedDataTypes = parsedDataTypes | ExtremeWind
		case "solar":
			parsedDataTypes = parsedDataTypes | Solar
		case "wind":
			parsedDataTypes = parsedDataTypes | Wind
		case "more_precip":
			parsedDataTypes = parsedDataTypes | MorePrecipitation
		case "more_phenomena":
			parsedDataTypes = parsedDataTypes | MorePhenomena
		case "soil_temperature":
			parsedDataTypes = parsedDataTypes | SoilTemperature
		case "water_equiv":
			parsedDataTypes = parsedDataTypes | WaterEquivalent
		case "weather_phenomena":
			parsedDataTypes = parsedDataTypes | WeatherPhenomena
		case "cloud_type":
			parsedDataTypes = parsedDataTypes | CloudType
		case "cloudiness":
			parsedDataTypes = parsedDataTypes | Cloudiness
		case "dew_point":
			parsedDataTypes = parsedDataTypes | DewPoint
		case "moisture":
			parsedDataTypes = parsedDataTypes | Moisture
		case "pressure":
			parsedDataTypes = parsedDataTypes | Pressure
		case "sun":
			parsedDataTypes = parsedDataTypes | Sun
		case "visibility":
			parsedDataTypes = parsedDataTypes | Visibility
		case "wind_synop":
			parsedDataTypes = parsedDataTypes | WindSynopsis
		case "climate_indices":
			parsedDataTypes = parsedDataTypes | ClimateIndices
		case "kl":
			parsedDataTypes = parsedDataTypes | StationObservations
		case "standard_format":
			parsedDataTypes = parsedDataTypes | StandardFormat
		case "wind_test":
			parsedDataTypes = parsedDataTypes | WindTest
		case "soil":
			parsedDataTypes = parsedDataTypes | Soil
		case "more_weather_phenomena":
			parsedDataTypes = parsedDataTypes | MoreWeatherPhenomena
		default:
			continue
		}
	}
	*dt = parsedDataTypes
}

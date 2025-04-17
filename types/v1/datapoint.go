package v1

import (
	"encoding/json"
	"strings"
)

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

// function.
func (dt *DataType) ParseString(s string) {
	dt.ParseStringWithSeperator(s, ",")
}

// function as it expects this by default.
func (dt *DataType) ParseStringWithSeperator(s, sep string) {
	var parsedDataTypes DataType
	parts := strings.Split(s, sep)
	for _, part := range parts {
		switch part {
		case "precipitation":
			parsedDataTypes |= Precipitation
		case "air_temperature":
			parsedDataTypes |= AirTemperature
		case "extreme_temperature":
			parsedDataTypes |= ExtremeTemperature
		case "extreme_wind":
			parsedDataTypes |= ExtremeWind
		case "solar":
			parsedDataTypes |= Solar
		case "wind":
			parsedDataTypes |= Wind
		case "more_precip":
			parsedDataTypes |= MorePrecipitation
		case "more_phenomena":
			parsedDataTypes |= MorePhenomena
		case "soil_temperature":
			parsedDataTypes |= SoilTemperature
		case "water_equiv":
			parsedDataTypes |= WaterEquivalent
		case "weather_phenomena":
			parsedDataTypes |= WeatherPhenomena
		case "cloud_type":
			parsedDataTypes |= CloudType
		case "cloudiness":
			parsedDataTypes |= Cloudiness
		case "dew_point":
			parsedDataTypes |= DewPoint
		case "moisture":
			parsedDataTypes |= Moisture
		case "pressure":
			parsedDataTypes |= Pressure
		case "sun":
			parsedDataTypes |= Sun
		case "visibility":
			parsedDataTypes |= Visibility
		case "wind_synop":
			parsedDataTypes |= WindSynopsis
		case "climate_indices":
			parsedDataTypes |= ClimateIndices
		case "kl":
			parsedDataTypes |= StationObservations
		case "standard_format":
			parsedDataTypes |= StandardFormat
		case "wind_test":
			parsedDataTypes |= WindTest
		case "soil":
			parsedDataTypes |= Soil
		case "more_weather_phenomena":
			parsedDataTypes |= MoreWeatherPhenomena
		default:
			continue
		}
	}
	*dt = parsedDataTypes
}

func (dt DataType) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.String())
}

func (dt *DataType) UnmarshalJSON(src []byte) error {
	var dataType string
	if err := json.Unmarshal(src, &dataType); err != nil {
		return err
	}
	dt.ParseString(dataType)
	return nil
}

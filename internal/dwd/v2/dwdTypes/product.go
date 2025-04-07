package dwdTypes

import (
	"encoding/json"
	"errors"
)

// Product represents a data product offered on the OpenData Portal.
type Product uint

const (
	ClimateObservation_AirTemperature Product = iota
	ClimateObservation_ClimateIndices
	ClimateObservation_Cloudiness
	ClimateObservation_CloudType
	ClimateObservation_DewPoint
	ClimateObservation_ExtremeTemperature
	ClimateObservation_ExtremeWind
	ClimateObservation_Moisture
	ClimateObservation_MorePhenomena
	ClimateObservation_MorePrecipitation
	ClimateObservation_MoreWeatherPhenomena
	ClimateObservation_Precipitation
	ClimateObservation_Pressure
	ClimateObservation_Soil
	ClimateObservation_SoilTemperature
	ClimateObservation_SolarRadiation
	ClimateObservation_Sun
	ClimateObservation_Visibility
	ClimateObservation_WaterEquivalent
	ClimateObservation_WeatherPhenomena
	ClimateObservation_WindSpeeds
	ClimateObservation_WindSynopsis
)

func (p Product) String() string {
	switch p {
	case ClimateObservation_AirTemperature:
		return "airTemperature"
	case ClimateObservation_ClimateIndices:
		return "climateIndices"
	case ClimateObservation_Cloudiness:
		return "cloudiness"
	case ClimateObservation_CloudType:
		return "cloudType"
	case ClimateObservation_DewPoint:
		return "dewPoint"
	case ClimateObservation_ExtremeTemperature:
		return "extremeTemperatures"
	case ClimateObservation_ExtremeWind:
		return "extremeWinds"
	case ClimateObservation_Moisture:
		return "moisture"
	case ClimateObservation_MorePhenomena:
		return "morePhenomena"
	case ClimateObservation_MorePrecipitation:
		return "morePrecipitation"
	case ClimateObservation_MoreWeatherPhenomena:
		return "moreWeatherPhenomena"
	case ClimateObservation_Precipitation:
		return "precipitation"
	case ClimateObservation_Pressure:
		return "pressure"
	case ClimateObservation_Soil:
		return "soil"
	case ClimateObservation_SoilTemperature:
		return "soilTemperature"
	case ClimateObservation_SolarRadiation:
		return "solarRadiation"
	case ClimateObservation_Sun:
		return "sun"
	case ClimateObservation_Visibility:
		return "visibility"
	case ClimateObservation_WaterEquivalent:
		return "waterEquivalent"
	case ClimateObservation_WeatherPhenomena:
		return "weatherPhenomena"
	case ClimateObservation_WindSpeeds:
		return "windSpeeds"
	case ClimateObservation_WindSynopsis:
		return "windSynopsis"
	default:
		return ""
	}
}

func (p Product) UrlPart() string {
	switch p {
	case ClimateObservation_AirTemperature:
		return "air_temperature"
	case ClimateObservation_ClimateIndices:
		return "climate_indices"
	case ClimateObservation_CloudType:
		return "cloud_type"
	case ClimateObservation_DewPoint:
		return "dew_point"
	case ClimateObservation_ExtremeTemperature:
		return "extreme_temperature"
	case ClimateObservation_ExtremeWind:
		return "extreme_wind"
	case ClimateObservation_MorePhenomena:
		return "more_phenomena"
	case ClimateObservation_MorePrecipitation:
		return "more_precip"
	case ClimateObservation_MoreWeatherPhenomena:
		return "more_weather_phenomena"
	case ClimateObservation_SoilTemperature:
		return "soil_temperature"
	case ClimateObservation_SolarRadiation:
		return "solar"
	case ClimateObservation_WaterEquivalent:
		return "water_equiv"
	case ClimateObservation_WeatherPhenomena:
		return "weather_phenomena"
	case ClimateObservation_WindSpeeds:
		return "wind"
	case ClimateObservation_WindSynopsis:
		return "wind_synop"
	default:
		return p.String()
	}
}

func (p *Product) Parse(src any) error {
	var productString string
	switch src := src.(type) {
	case []byte:
		productString = string(src)
	case string:
		productString = src
	case uint:
		*p = Product(src)
		return nil
	default:
		return errors.New("unsupported input type")
	}

	switch productString {
	case ClimateObservation_AirTemperature.String(), ClimateObservation_AirTemperature.UrlPart():
		*p = ClimateObservation_AirTemperature
	case ClimateObservation_ClimateIndices.String(), ClimateObservation_ClimateIndices.UrlPart():
		*p = ClimateObservation_ClimateIndices
	case ClimateObservation_Cloudiness.String(), ClimateObservation_Cloudiness.UrlPart():
		*p = ClimateObservation_Cloudiness
	case ClimateObservation_CloudType.String(), ClimateObservation_CloudType.UrlPart():
		*p = ClimateObservation_CloudType
	case ClimateObservation_DewPoint.String(), ClimateObservation_DewPoint.UrlPart():
		*p = ClimateObservation_DewPoint
	case ClimateObservation_ExtremeTemperature.String(), ClimateObservation_ExtremeTemperature.UrlPart():
		*p = ClimateObservation_ExtremeTemperature
	case ClimateObservation_ExtremeWind.String(), ClimateObservation_ExtremeWind.UrlPart():
		*p = ClimateObservation_ExtremeWind
	case ClimateObservation_Moisture.String(), ClimateObservation_Moisture.UrlPart():
		*p = ClimateObservation_Moisture
	case ClimateObservation_MorePhenomena.String(), ClimateObservation_MorePhenomena.UrlPart():
		*p = ClimateObservation_MorePhenomena
	case ClimateObservation_MorePrecipitation.String(), ClimateObservation_MorePrecipitation.UrlPart():
		*p = ClimateObservation_MorePrecipitation
	case ClimateObservation_MoreWeatherPhenomena.String(), ClimateObservation_MoreWeatherPhenomena.UrlPart():
		*p = ClimateObservation_MoreWeatherPhenomena
	case ClimateObservation_Precipitation.String(), ClimateObservation_Precipitation.UrlPart():
		*p = ClimateObservation_Precipitation
	case ClimateObservation_Pressure.String(), ClimateObservation_Pressure.UrlPart():
		*p = ClimateObservation_Pressure
	case ClimateObservation_Soil.String(), ClimateObservation_Soil.UrlPart():
		*p = ClimateObservation_Soil
	case ClimateObservation_SoilTemperature.String(), ClimateObservation_SoilTemperature.UrlPart():
		*p = ClimateObservation_SoilTemperature
	case ClimateObservation_SolarRadiation.String(), ClimateObservation_SolarRadiation.UrlPart():
		*p = ClimateObservation_SolarRadiation
	case ClimateObservation_Sun.String(), ClimateObservation_Sun.UrlPart():
		*p = ClimateObservation_Sun
	case ClimateObservation_Visibility.String(), ClimateObservation_Visibility.UrlPart():
		*p = ClimateObservation_Visibility
	case ClimateObservation_WaterEquivalent.String(), ClimateObservation_WaterEquivalent.UrlPart():
		*p = ClimateObservation_WaterEquivalent
	case ClimateObservation_WeatherPhenomena.String(), ClimateObservation_WeatherPhenomena.UrlPart():
		*p = ClimateObservation_WeatherPhenomena
	case ClimateObservation_WindSpeeds.String(), ClimateObservation_WindSpeeds.UrlPart():
		*p = ClimateObservation_WindSpeeds
	case ClimateObservation_WindSynopsis.String(), ClimateObservation_WindSynopsis.UrlPart():
		*p = ClimateObservation_WindSynopsis
	default:
		return errors.New("unsupported product")
	}
	return nil
}

func (p Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Product) UnmarshalJSON(src []byte) error {
	return p.Parse(src)
}

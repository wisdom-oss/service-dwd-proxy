package v2

import "microservice/internal/dwd/v2/dwdTypes"

type Granularity = dwdTypes.Granularity
type Product = dwdTypes.Product

// AvailableClimateObservationProducts contains a mapping of the products to the
// each of the available granularities.
var AvailableClimateObservationProducts = map[Granularity][]Product{
	dwdTypes.Granularity_EveryMinute: {dwdTypes.ClimateObservation_Precipitation},
	dwdTypes.Granularity_Every5Mins:  {dwdTypes.ClimateObservation_Precipitation},
	dwdTypes.Granularity_Every10Mins: {
		dwdTypes.ClimateObservation_AirTemperature,
		dwdTypes.ClimateObservation_ExtremeTemperature,
		dwdTypes.ClimateObservation_ExtremeWind,
		dwdTypes.ClimateObservation_Precipitation,
		dwdTypes.ClimateObservation_SolarRadiation,
		dwdTypes.ClimateObservation_WindSpeeds,
	},
	dwdTypes.Granularity_Hourly: {
		dwdTypes.ClimateObservation_AirTemperature,
		dwdTypes.ClimateObservation_CloudType,
		dwdTypes.ClimateObservation_Cloudiness,
		dwdTypes.ClimateObservation_DewPoint,
		dwdTypes.ClimateObservation_ExtremeWind,
		dwdTypes.ClimateObservation_Moisture,
		dwdTypes.ClimateObservation_Precipitation,
		dwdTypes.ClimateObservation_Pressure,
		dwdTypes.ClimateObservation_SoilTemperature,
		dwdTypes.ClimateObservation_SolarRadiation,
		dwdTypes.ClimateObservation_Sun,
		dwdTypes.ClimateObservation_Visibility,
		dwdTypes.ClimateObservation_WeatherPhenomena,
		dwdTypes.ClimateObservation_WindSpeeds,
		dwdTypes.ClimateObservation_WindSynopsis,
	},
	dwdTypes.Granularity_SubDaily: {
		dwdTypes.ClimateObservation_AirTemperature,
		dwdTypes.ClimateObservation_Cloudiness,
		dwdTypes.ClimateObservation_ExtremeWind,
		dwdTypes.ClimateObservation_Moisture,
		dwdTypes.ClimateObservation_Pressure,
		dwdTypes.ClimateObservation_Soil,
		dwdTypes.ClimateObservation_Visibility,
		dwdTypes.ClimateObservation_WindSpeeds,
	},
	dwdTypes.Granularity_Daily: {
		dwdTypes.ClimateObservation_MorePrecipitation,
		dwdTypes.ClimateObservation_MoreWeatherPhenomena,
		dwdTypes.ClimateObservation_SoilTemperature,
		dwdTypes.ClimateObservation_SolarRadiation,
		dwdTypes.ClimateObservation_WaterEquivalent,
		dwdTypes.ClimateObservation_WeatherPhenomena,
	},
	dwdTypes.Granularity_Monthly: {
		dwdTypes.ClimateObservation_ClimateIndices,
		dwdTypes.ClimateObservation_MorePrecipitation,
		dwdTypes.ClimateObservation_WeatherPhenomena,
	},
	dwdTypes.Granularity_Annual: {
		dwdTypes.ClimateObservation_ClimateIndices,
		dwdTypes.ClimateObservation_MorePrecipitation,
		dwdTypes.ClimateObservation_WeatherPhenomena,
	},
}

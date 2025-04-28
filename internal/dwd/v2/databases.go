package v2

const (
	ClimateObservationsUrlKey  = "climateObservations"
	ClimateObservationsBaseUrl = "https://opendata.dwd.de/climate_environment/CDC/observations_germany/climate/"
)

var Databases = map[string]string{
	ClimateObservationsUrlKey: ClimateObservationsBaseUrl,
}

var Products = map[string]map[Granularity][]Product{
	ClimateObservationsUrlKey: AvailableClimateObservationProducts,
}

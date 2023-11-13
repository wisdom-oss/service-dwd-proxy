package routes

import "net/http"

// BaseHost contains the opendata portals base host
const BaseHost = "https://opendata.dwd.de"

// BasePath contains the path to the climate observations
const BasePath = "/climate_environment/CDC/observations_germany/climate"

// client is the HTTP Client used to access the open data portal
var client = http.Client{}

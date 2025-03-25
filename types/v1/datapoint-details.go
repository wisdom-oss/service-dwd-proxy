package v1

import "time"

type DatapointDetails struct {
	Resolution Resolution `json:"resolution"`
	From       time.Time  `json:"availableFrom"`
	Until      time.Time  `json:"availableUntil"`
}

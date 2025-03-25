package v1

import "time"

type Capability struct {
	// Capability sets the capability that is described
	DataType       DataType   `json:"dataType"`
	Resolution     Resolution `json:"resolution"`
	AvailableFrom  time.Time  `json:"availableFrom"`
	AvailableUntil time.Time  `json:"availableUntil"`
}

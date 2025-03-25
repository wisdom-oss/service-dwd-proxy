package v1

import (
	"time"
)

// been recorded.
type TimeseriesField struct {
	// Name contains the name which identifies this field in a timeseries
	Name string `json:"name"`
	// Description contains a textual description of the field. It mainly is
	// written in German since the DWD only supplies a German description for
	// the fields
	Description string `json:"description"`
	// AvailableFrom denotes the date at which the field has been recorded for
	// the first time
	AvailableFrom time.Time `json:"availableFrom"`
	// AvailableUntil denotes the date at which the field as been recorded for
	// the last time
	AvailableUntil time.Time `json:"availableUntil"`
	// Unit denotes the unit used for the recorded measurements
	Unit string `json:"unit"`
}

package v2

import "time"

type DateTimeRange struct {
	Start time.Time `json:"from"`
	End   time.Time `json:"until"`
}

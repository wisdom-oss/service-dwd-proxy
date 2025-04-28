package v2

import (
	"time"

	"microservice/internal/dwd/v2/dwdTypes"
)

type Datapoint struct {
	Label        string                `json:"label"`
	Timestamp    time.Time             `json:"ts"`
	Value        any                   `json:"value"`
	Unit         *string               `json:"unit"`
	QualityLevel *dwdTypes.QualityFlag `json:"qualityLevel"`
}

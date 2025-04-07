package v2

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"

	"microservice/internal/dwd/v2/dwdTypes"
)

type Station struct {
	ID                string                            `csv:"id"`
	Name              string                            `csv:"name"`
	Height            float64                           `csv:"height"`
	Location          *geom.Point                       `csv:"geometry"`
	SupportedProducts map[string][]dwdTypes.Granularity `csv:"-"`
}

func (s Station) MarshalJSON() ([]byte, error) {

	return s.ToFeature().MarshalJSON()
}

func (s Station) ToFeature() *geojson.Feature {

	f := geojson.Feature{
		ID:       s.ID,
		Geometry: s.Location,
		Properties: map[string]any{
			"name":     s.Name,
			"products": s.SupportedProducts,
		},
	}
	return &f
}

package types

import "github.com/paulmach/go.geojson"

type Station struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	State    string           `json:"state"`
	Height   float64          `json:"height"`
	Location geojson.Geometry `json:"location"`
}

// SetLocation updates the location of the current station to one with the
// coordinates supplied to the function
func (s *Station) SetLocation(longitude, latitude float64) {
	geom := geojson.NewPointGeometry([]float64{longitude, latitude})
	s.Location = *geom
}

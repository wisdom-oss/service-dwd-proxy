package types

import "github.com/paulmach/go.geojson"

type Station struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	State            string                `json:"state"`
	Height           float64               `json:"height"`
	Location         *geojson.Geometry     `json:"location"`
	DataCapabilities map[string]Resolution `json:"capabilities"`
}

func (s *Station) UpdateCapabilities(dataType DataType, resolution Resolution) {
	if s.DataCapabilities == nil {
		s.DataCapabilities = make(map[string]Resolution)
	}
	supportedResolutions, resolutionsAvailable := s.DataCapabilities[dataType.String()]
	if resolutionsAvailable {
		supportedResolutions = supportedResolutions | resolution
		s.DataCapabilities[dataType.String()] = supportedResolutions
	} else {
		s.DataCapabilities[dataType.String()] = resolution
	}
}

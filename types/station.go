package types

import (
	"slices"
	"time"

	"github.com/paulmach/go.geojson"
)

type Station struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	State            string            `json:"state"`
	Height           float64           `json:"height"`
	Location         *geojson.Geometry `json:"location"`
	Historical       bool              `json:"historical"`
	DataCapabilities []Capability      `json:"capabilities"`
}

func (s *Station) AddCapability(c Capability) {
	// check if the capability is already available in the capability array
	capabilityKnown := false
	capabilityIndex := -1
	for index, knownCapability := range s.DataCapabilities {
		if knownCapability.DataType == c.DataType && knownCapability.Resolution == c.Resolution {
			capabilityKnown = true
			capabilityIndex = index
			break
		}
	}
	if capabilityKnown {
		// now get the known capability
		knownCapability := s.DataCapabilities[capabilityIndex]
		// now remove this capability from the array to allow editing it
		capabilityArray := slices.Delete(s.DataCapabilities, capabilityIndex, capabilityIndex)
		// now update the known capability if needed
		if knownCapability.AvailableFrom.After(c.AvailableFrom) {
			knownCapability.AvailableFrom = c.AvailableFrom
		}
		if knownCapability.AvailableUntil.Before(c.AvailableUntil) {
			knownCapability.AvailableUntil = c.AvailableUntil
		}
		// now push the capability in the array again
		capabilityArray = append(capabilityArray, knownCapability)
		s.DataCapabilities = capabilityArray
		return
	}
	s.DataCapabilities = append(s.DataCapabilities, c)
}

// UpdateHistoricalState iterates over all data capabilities and checks if
// at least one capability has set the current date as last update date for
// a capability.
// If this is not the case, the station's historical state is activated
func (s *Station) UpdateHistoricalState() error {
	var currentDayCapabilities bool
	// construct today's date
	today, err := time.Parse("20060102", time.Now().Format("20060102"))
	if err != nil {
		return err
	}

	for _, c := range s.DataCapabilities {
		if c.AvailableUntil.Equal(today) || c.AvailableFrom.Equal(today) {
			currentDayCapabilities = true
			break
		}
	}

	s.Historical = !currentDayCapabilities
	return nil
}

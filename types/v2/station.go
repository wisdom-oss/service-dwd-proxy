package v2

import (
	"slices"
	"time"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

type Station struct {
	ID                string                             `csv:"id"`
	Name              string                             `csv:"name"`
	Height            float64                            `csv:"height"`
	Location          *geom.Point                        `csv:"geometry"`
	SupportedProducts map[string]map[string][2]time.Time `csv:"-"`
}

func (s Station) MarshalJSON() ([]byte, error) {

	return s.ToFeature().MarshalJSON()
}

func (s Station) ToFeature() *geojson.Feature {

	products := make(map[string]map[string]struct {
		From  time.Time `json:"from"`
		Until time.Time `json:"until"`
	})

	for product, granularityAvailablility := range s.SupportedProducts {
		granularities := make([]string, 0, len(s.SupportedProducts[product]))
		for granularity := range s.SupportedProducts[product] {
			granularities = append(granularities, granularity)
		}

		for _, granulartiy := range granularities {
			if products[product] == nil {
				products[product] = make(map[string]struct {
					From  time.Time `json:"from"`
					Until time.Time `json:"until"`
				})
			}
			products[product][granulartiy] = struct {
				From  time.Time `json:"from"`
				Until time.Time `json:"until"`
			}{
				From:  granularityAvailablility[granulartiy][0],
				Until: granularityAvailablility[granulartiy][1],
			}

		}
	}

	f := geojson.Feature{
		ID:       s.ID,
		Geometry: s.Location,
		Properties: map[string]any{
			"name":     s.Name,
			"products": products,
			"id":       s.ID,
		},
	}
	return &f
}

func (this *Station) MergeProducts(other Station) {
	for otherProduct, granularityAvaiability := range other.SupportedProducts {
		_, found := this.SupportedProducts[otherProduct]
		if !found {
			this.SupportedProducts[otherProduct] = granularityAvaiability
			continue
		}

		thisGranularities := make([]string, 0, len(this.SupportedProducts[otherProduct]))
		for granularity := range this.SupportedProducts[otherProduct] {
			thisGranularities = append(thisGranularities, granularity)
		}

		otherGranularities := make([]string, 0, len(other.SupportedProducts[otherProduct]))
		for granularity := range other.SupportedProducts[otherProduct] {
			otherGranularities = append(otherGranularities, granularity)
		}

		for _, granularity := range otherGranularities {
			if !slices.Contains(thisGranularities, granularity) {
				this.SupportedProducts[otherProduct][granularity] = other.SupportedProducts[otherProduct][granularity]
				continue
			}

			thisStart := this.SupportedProducts[otherProduct][granularity][0]
			thisEnd := this.SupportedProducts[otherProduct][granularity][1]

			otherStart := other.SupportedProducts[otherProduct][granularity][0]
			otherEnd := other.SupportedProducts[otherProduct][granularity][1]

			var newTimes [2]time.Time

			if otherStart.Before(thisStart) {
				newTimes[0] = otherStart
			} else {
				newTimes[0] = thisStart
			}

			if otherEnd.After(thisEnd) {
				newTimes[1] = otherEnd
			} else {
				newTimes[1] = thisEnd
			}

			other.SupportedProducts[otherProduct][granularity] = newTimes
		}

	}
}

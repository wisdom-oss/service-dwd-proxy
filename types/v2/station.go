package v2

import (
	"slices"
	"time"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"

	"microservice/internal/dwd/v2/dwdTypes"
)

type Station struct {
	ID                string                                                      `csv:"id"`
	Name              string                                                      `csv:"name"`
	Height            float64                                                     `csv:"height"`
	Location          *geom.Point                                                 `csv:"geometry"`
	SupportedProducts map[dwdTypes.Product]map[dwdTypes.Granularity]DateTimeRange `csv:"-"`
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
		granularities := make([]dwdTypes.Granularity, 0, len(s.SupportedProducts[product]))
		for granularity := range s.SupportedProducts[product] {
			granularities = append(granularities, granularity)
		}

		for _, granulartiy := range granularities {
			if products[product.String()] == nil {
				products[product.String()] = make(map[string]struct {
					From  time.Time `json:"from"`
					Until time.Time `json:"until"`
				})
			}
			products[product.String()][granulartiy.String()] = struct {
				From  time.Time `json:"from"`
				Until time.Time `json:"until"`
			}{
				From:  granularityAvailablility[granulartiy].Start,
				Until: granularityAvailablility[granulartiy].End,
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

		thisGranularities := make([]dwdTypes.Granularity, 0, len(this.SupportedProducts[otherProduct]))
		for granularity := range this.SupportedProducts[otherProduct] {
			thisGranularities = append(thisGranularities, granularity)
		}

		otherGranularities := make([]dwdTypes.Granularity, 0, len(other.SupportedProducts[otherProduct]))
		for granularity := range other.SupportedProducts[otherProduct] {
			otherGranularities = append(otherGranularities, granularity)
		}

		for _, granularity := range otherGranularities {
			if !slices.Contains(thisGranularities, granularity) {
				this.SupportedProducts[otherProduct][granularity] = other.SupportedProducts[otherProduct][granularity]
				continue
			}

			thisStart := this.SupportedProducts[otherProduct][granularity].Start
			thisEnd := this.SupportedProducts[otherProduct][granularity].End

			otherStart := other.SupportedProducts[otherProduct][granularity].Start
			otherEnd := other.SupportedProducts[otherProduct][granularity].End

			var r DateTimeRange

			if otherStart.Before(thisStart) {
				r.Start = otherStart
			} else {
				r.Start = thisStart
			}

			if otherEnd.After(thisEnd) {
				r.End = otherEnd
			} else {
				r.End = thisEnd
			}

			other.SupportedProducts[otherProduct][granularity] = r
		}

	}
}

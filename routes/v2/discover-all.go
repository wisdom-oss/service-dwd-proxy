package v2

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/twpayne/go-geom/encoding/geojson"
	"golang.org/x/sync/errgroup"

	dwd "microservice/internal/dwd/v2"
	v2 "microservice/types/v2"
)

func DiscoverAllStations(c *gin.Context) {
	var paralel errgroup.Group
	var arrayLock sync.Mutex
	var allStations []v2.Station

	for granularity, products := range dwd.AvailableClimateObservationProducts {
		for _, product := range products {
			paralel.Go(func() error {
				discoveredStations, err := dwd.DiscoverStations(dwd.ClimateObservationsBaseUrl, granularity, product)
				if err != nil {
					return err
				}

				for _, station := range discoveredStations {
					arrayLock.Lock()
					allStations = append(allStations, station)
					arrayLock.Unlock()
				}
				return nil
			})

		}
	}

	err := paralel.Wait()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	mergedStations := make(map[string]v2.Station)

	for _, station := range allStations {
		mapKey := fmt.Sprintf("%f|%f", station.Location.X(), station.Location.Y())

		processedStation, alreadyProcessed := mergedStations[mapKey]
		if !alreadyProcessed {
			mergedStations[mapKey] = station
			continue
		}

		processedStation.MergeProducts(station)

		mergedStations[mapKey] = processedStation
	}

	var features []*geojson.Feature //nolint:prealloc
	for _, station := range mergedStations {
		features = append(features, station.ToFeature())
	}

	featureCollection := geojson.FeatureCollection{Features: features}

	c.JSON(http.StatusOK, &featureCollection)

}

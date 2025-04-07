package v2

import (
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twpayne/go-geom/encoding/geojson"
	"golang.org/x/sync/errgroup"

	dwd "microservice/internal/dwd/v2"
	"microservice/internal/dwd/v2/dwdTypes"
	v2 "microservice/types/v2"
)

func DiscoverAllStations(c *gin.Context) {
	var paralel errgroup.Group
	var arrayLock sync.Mutex
	var allStations []v2.Station

	handlingStart := time.Now()

	for granularity, products := range dwd.AvailableClimateObservationProducts {
		for _, product := range products {
			paralel.Go(func() error {
				discoveredStations, err := dwd.DiscoverStations(dwd.ClimateObservationsBaseUrl, granularity, product)
				if err != nil {
					return err
				}

				for _, discoveredStation := range discoveredStations {
					discoveredStation.SupportedProducts = map[string][]dwdTypes.Granularity{
						product.String(): {granularity},
					}
					arrayLock.Lock()
					allStations = append(allStations, discoveredStation)
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

	stationDiscoveryDuration := time.Since(handlingStart)
	stationDiscoveryMillis := float64(stationDiscoveryDuration.Microseconds()) / 1000

	mergedStations := make(map[string]v2.Station)

	for _, station := range allStations {
		mapKey := fmt.Sprintf("%f|%f", station.Location.X(), station.Location.Y())

		processedStation, alreadyProcessed := mergedStations[mapKey]
		if !alreadyProcessed {
			mergedStations[mapKey] = station
			continue
		}

		for product, granularity := range station.SupportedProducts {
			if slices.Contains(processedStation.SupportedProducts[product], granularity[0]) {
				continue
			}
			processedStation.SupportedProducts[product] = append(processedStation.SupportedProducts[product], granularity...)
		}

		mergedStations[mapKey] = processedStation
	}

	stationMergingDuration := time.Since(handlingStart) - stationDiscoveryDuration
	stationMergingMillis := float64(stationMergingDuration.Microseconds()) / 1000

	var features []*geojson.Feature //nolint:prealloc
	for _, station := range mergedStations {
		features = append(features, station.ToFeature())
	}

	featureCollection := geojson.FeatureCollection{Features: features}

	timingHeaderValue := `discovery;dur=%f;desc="Station Discovery and File Parsing", merging;dur=%f;desc="Station Merging"`
	c.Header("Server-Timing", fmt.Sprintf(timingHeaderValue, stationDiscoveryMillis, stationMergingMillis))
	c.JSON(http.StatusOK, &featureCollection)

}

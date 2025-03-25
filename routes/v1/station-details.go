package v1

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"

	"microservice/internal/redis"
	types "microservice/types/v1"
)

func StationDetails(c *gin.Context) {
	client := redis.Client()

	// retrieve the stored station list from redis
	stationListBytes, err := client.Get(c, RedisKey_StationList).Bytes()
	if err != nil {
		c.Abort()
		if redis.IsNotFound(err) {
			errRedisCacheUnprimed.Emit(c)
			return
		}
		_ = c.Error(err)
		return
	}

	// create a brotli reader to allow data decompression
	reader := brotli.NewReader(bytes.NewReader(stationListBytes))

	var stations []types.Station
	err = json.NewDecoder(reader).Decode(&stations)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	stationID := c.Param("stationID")
	for _, station := range stations {
		if station.ID != stationID {
			continue
		}

		c.JSON(http.StatusOK, station)
		return
	}

	c.Abort()
	errUnknownStation.Emit(c)
}

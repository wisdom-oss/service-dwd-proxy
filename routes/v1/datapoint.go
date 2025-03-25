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

func DatapointInformation(c *gin.Context) {
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

	datapoint := types.DataType(0)
	datapoint.ParseString(c.Param("datapoint"))

	if datapoint == 0 {
		errUnsupportedDatapoint.Emit(c)
		return
	}

	var station types.Station

	for idx, s := range stations {
		if s.ID == c.Param("stationID") {
			station = s
			break
		}
		if idx == len(stations)-1 {
			c.Abort()
			errUnknownStation.Emit(c)
			return
		}
	}

	datapointDetails := []types.DatapointDetails{}
	for _, capability := range station.DataCapabilities {
		if capability.DataType == datapoint {
			datapointDetails = append(datapointDetails, types.DatapointDetails{
				Resolution: capability.Resolution,
				From:       capability.AvailableFrom,
				Until:      capability.AvailableUntil,
			})
		}
	}

	c.JSON(http.StatusOK, datapointDetails)
}

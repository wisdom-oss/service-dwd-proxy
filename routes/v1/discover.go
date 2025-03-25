package v1

import (
	"bytes"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"

	"microservice/internal/redis"
)

// Discover reads from the primed redis cache and returns the contents directly.
func Discover(c *gin.Context) {
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

	// respond with the data contained in the station list
	c.DataFromReader(http.StatusOK, -1, contentType, reader, nil)
}

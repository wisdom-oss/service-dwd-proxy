package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dwd "microservice/internal/dwd/v2"
	v2 "microservice/types/v2"
)

func ValidateConnection(c *gin.Context) {
	health := make(map[string]v2.HealthStatus)
	for name, url := range dwd.Databases {
		res, err := http.Get(url) //nolint:gosec // The variable urls are from our own constants
		if err != nil {
			health[name] = v2.HealthStatus{Healthy: false, Reason: err.Error()}
			continue
		}

		if res.StatusCode != http.StatusOK {
			health[name] = v2.HealthStatus{Healthy: false, Reason: "response code indicated not ok"}
			continue
		}

		health[name] = v2.HealthStatus{Healthy: true}
	}
	c.JSON(http.StatusOK, health)
}

package v1

import (
	"net/http"

	"github.com/wisdom-oss/common-go/v3/types"
)

var errRedisCacheUnprimed = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.6.4",
	Status: http.StatusServiceUnavailable,
	Title:  "Redis Unprimed",
	Detail: "The redis database currently does not contain the required data. Please try again later",
}

var errUnknownStation = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
	Status: http.StatusNotFound,
	Title:  "Unknown Station",
	Detail: "The station is not known.",
}

var errUnsupportedDatapoint = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Datapoint Unsupported",
	Detail: "The specified datapoint is not supported by the station or is invalid",
}

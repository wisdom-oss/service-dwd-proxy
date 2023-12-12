package routes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"

	"github.com/wisdom-oss/service-dwd-proxy/globals"
)

// DiscoverMetadata makes a request to the OpenData portal hosted by the
// DWD.
// It returns a response marking every data type available and the time
// resolutions available for it.
// Furthermore, it returns the stations with their available data points
// and time resolutions.
func DiscoverMetadata(w http.ResponseWriter, r *http.Request) {
	// get the error handler
	errorHandler := r.Context().Value("error-channel").(chan<- interface{})
	statusChannel := r.Context().Value("status-channel").(<-chan bool)

	stationListBytes, err := globals.RedisClient.Get(r.Context(), "dwd-station-list").Bytes()
	if err != nil {
		errorHandler <- fmt.Errorf("unable to get station list from redis: %w", err)
		<-statusChannel
		return
	}

	// now read the bytes into a json response back from brotli
	byteReader := bytes.NewReader(stationListBytes)
	brotliReader := brotli.NewReader(byteReader)

	w.Header().Set("Content-Type", "application/json")
	_, err = io.Copy(w, brotliReader)
	if err != nil {
		errorHandler <- fmt.Errorf("unable return response: %w", err)
		<-statusChannel
		return
	}

}

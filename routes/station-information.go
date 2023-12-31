package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/go-chi/chi/v5"

	"github.com/wisdom-oss/service-dwd-proxy/globals"
	"github.com/wisdom-oss/service-dwd-proxy/types"
)

func StationInformation(w http.ResponseWriter, r *http.Request) {
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

	uncompressedBytes, err := io.ReadAll(brotliReader)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to decompress station list: %w", err)
		<-statusChannel
		return
	}
	// now parse the json array stored in the brotli reader
	var stations []types.Station
	err = json.Unmarshal(uncompressedBytes, &stations)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to parse uncompressed bytes: %w", err)
		<-statusChannel
		return
	}
	// now search for the station by iterating over them
	stationID := chi.URLParam(r, "stationID")

	for _, station := range stations {
		if station.ID == stationID {
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(station)
			return
		}
	}

	// since no station has been found return a 404
	errorHandler <- "UNKNOWN_STATION"
	<-statusChannel
	return
}

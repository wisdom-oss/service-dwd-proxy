package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/go-chi/chi/v5"

	"github.com/wisdom-oss/service-dwd-proxy/globals"
	"github.com/wisdom-oss/service-dwd-proxy/types"
)

func DataTypeInformation(w http.ResponseWriter, r *http.Request) {
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
	dataTypeString := chi.URLParam(r, "dataType")
	// now check if the data type is supported
	dataType := types.DataType(0)
	dataType.ParseString(dataTypeString)
	if dataType == 0 {
		errorHandler <- "UNSUPPORTED_DATA_TYPE"
		<-statusChannel
		return
	}
	var station types.Station
	var stationFound bool
	for _, station = range stations {
		if station.ID == stationID {
			stationFound = true
			break
		}
	}
	if !stationFound {
		// since no station has been found return a 404
		errorHandler <- "UNKNOWN_STATION"
		<-statusChannel
		return
	}

	type resolutionTimeRange struct {
		Resolution     types.Resolution `json:"resolution"`
		AvailableFrom  time.Time        `json:"availableFrom"`
		AvailableUntil time.Time        `json:"availableUntil"`
	}

	var supportedResolutions []resolutionTimeRange

	for _, capability := range station.DataCapabilities {
		if capability.DataType == dataType {
			supportedResolutions = append(supportedResolutions, resolutionTimeRange{
				Resolution:     capability.Resolution,
				AvailableFrom:  capability.AvailableFrom,
				AvailableUntil: capability.AvailableUntil,
			})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(supportedResolutions)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

}

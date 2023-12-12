package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/wisdom-oss/service-dwd-proxy/helpers"
	"github.com/wisdom-oss/service-dwd-proxy/types"
)

// TimeSeries pulls the specified timeseries from the opendata portal
func TimeSeries(w http.ResponseWriter, r *http.Request) {
	// get the error handler
	errorHandler := r.Context().Value("error-channel").(chan<- interface{})
	statusChannel := r.Context().Value("status-channel").(<-chan bool)

	// get the url parameters
	rawDataType := chi.URLParam(r, "data-type")
	rawResolution := chi.URLParam(r, "resolution")
	stationID := chi.URLParam(r, "station")

	// get the query parameters
	queryParameter := r.URL.Query()
	fromTimeStampString := strings.TrimSpace(queryParameter.Get("from"))
	untilTimeStampString := strings.TrimSpace(queryParameter.Get("until"))

	var dataStart, dataEnd time.Time

	if fromTimeStampString != "" {
		i, err := strconv.ParseInt(fromTimeStampString, 10, 64)
		if err != nil {
			errorHandler <- "INVALID_TIMESTAMP"
			<-statusChannel
			return
		}

		dataStart = time.Unix(i, 0)
	}

	if untilTimeStampString != "" {
		i, err := strconv.ParseInt(untilTimeStampString, 10, 64)
		if err != nil {
			errorHandler <- "INVALID_TIMESTAMP"
			<-statusChannel
			return
		}

		dataEnd = time.Unix(i, 0)
	}

	// now check if the data type is supported
	dataType := types.DataType(0)
	dataType.ParseString(rawDataType)
	if dataType == 0 {
		errorHandler <- "UNSUPPORTED_DATA_TYPE"
		<-statusChannel
		return
	}

	resolution := types.Resolution(0)
	resolution.ParseString(rawResolution)
	if resolution == 0 {
		errorHandler <- "UNSUPPORTED_DATA_RESOLUTION"
		<-statusChannel
		return
	}

	// now build the initial query url to get the index page
	url := fmt.Sprintf("%s/%s/%s/%s", BaseHost, BasePath, resolution, dataType)
	// now try to get the index page
	indexPage, err := helpers.GetIndexPage(url)
	if err != nil {
		if errors.Is(err, helpers.ErrStatusNotFound) {
			errorHandler <- "UNSUPPORTED_DATA_TYPE"
			<-statusChannel
			return
		}
		if errors.Is(err, helpers.ErrStatusNot200) {
			errorHandler <- "WRONG_STATUS_CODE"
			<-statusChannel
			return
		}
		errorHandler <- err
		<-statusChannel
		return
	}
	// now parse the index page for the recent/historic/now folders
	folders := helpers.FilterDocumentForFolders(indexPage)
	var dataFileUrls []string
	for _, folder := range folders {
		url := fmt.Sprintf("%s/%s", url, folder)
		filePage, err := helpers.GetIndexPage(url)
		if err != nil {
			if errors.Is(err, helpers.ErrStatusNotFound) {
				errorHandler <- "UNSUPPORTED_DATA_TYPE"
				<-statusChannel
				return
			}
			if errors.Is(err, helpers.ErrStatusNot200) {
				errorHandler <- "WRONG_STATUS_CODE"
				<-statusChannel
				return
			}
			errorHandler <- err
			<-statusChannel
			return
		}
		availableFiles := helpers.FilterDocumentForFiles(filePage)
		for _, availableFile := range availableFiles {
			if strings.Contains(availableFile, stationID) {
				dataFileUrls = append(dataFileUrls, fmt.Sprintf("%s%s", url, availableFile))
			}
		}
	}
	if dataFileUrls == nil {
		errorHandler <- "STATION_WITHOUT_DATA"
		<-statusChannel
		return
	}
	var data []map[string]interface{}
	var parameters []types.TimeseriesField
	for _, dataFileUrl := range dataFileUrls {
		file, err := helpers.DownloadFile(dataFileUrl)
		if err != nil {
			errorHandler <- err
			<-statusChannel
			return
		}
		if strings.Contains(dataFileUrl, "meta_data") {
			fieldDescriptors, err := helpers.ParseMetadataArchive(file.Name())
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
			parameters = append(parameters, fieldDescriptors...)
		} else {
			datasets, fieldDescriptors, err := helpers.ParseDataFile(file.Name(), [2]time.Time{dataStart, dataEnd})
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
			parameters = append(parameters, fieldDescriptors...)
			data = append(data, datasets...)
		}

	}
	w.Header().Set("Content-Type", "application/json")
	type res struct {
		Data []map[string]interface{} `json:"timeseries"`
		Meta []types.TimeseriesField  `json:"metadata"`
	}
	response := res{
		Data: data,
		Meta: parameters,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		errorHandler <- err
		<-statusChannel
		return
	}

}

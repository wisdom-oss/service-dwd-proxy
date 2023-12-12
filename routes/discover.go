package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	geojson "github.com/paulmach/go.geojson"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"

	"github.com/wisdom-oss/service-dwd-proxy/helpers"
	"github.com/wisdom-oss/service-dwd-proxy/types"
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

	// request the initial page of the index for the climate observations
	url := fmt.Sprintf("%s%s", BaseHost, BasePath)
	res, err := client.Get(url)
	if err != nil {
		errorHandler <- fmt.Errorf("unable to access index page of open data portal: %w", err)
		<-statusChannel
		return
	}
	if res.StatusCode != 200 {
		log.Error().Int("httpStatusCode", res.StatusCode).Msg("unexpected http code in response")
		errorHandler <- "WRONG_STATUS_CODE"
		<-statusChannel
		return
	}
	// now parse the html tokens in on the page
	document, err := html.Parse(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse html page of dwd index page")
		errorHandler <- fmt.Errorf("unable to parse the index html page: %w", err)
		<-statusChannel
		return
	}
	// now resolve the folders that are listed in the document to be able to
	// find the currently available data resolutionFolders
	resolutionFolders := helpers.FilterDocumentForFolders(document)

	// to allow an asynchronous parsing of all resolutions and their available
	// data points, create a channel which will collect the results returned
	// by the http calls and folder parsing.
	// furthermore, create a wait group to allow the resynchronisation after
	// all needed discoveries are done.
	// create a boolean as well, indicating if an error happened to let the
	// handler exit properly
	var waitGroup sync.WaitGroup
	var errorOccurred bool
	// since the following actions are based on the number of detected
	// resolution folders, add the number of these to the wait group
	waitGroup.Add(len(resolutionFolders))

	// now iterate over the resolution folders and find the different data types
	dataTypeUrlChan := make(chan []string, len(resolutionFolders))
	for _, resolutionFolder := range resolutionFolders {
		resolutionFolder := resolutionFolder
		go func() {
			// wait until this function exits to notify the wait group of this
			defer waitGroup.Done()
			// now build the url
			url := fmt.Sprintf("%s/%s", url, resolutionFolder)
			// now request the index page
			log.Debug().Str("url", url).Msg("requesting index page")
			page, err := helpers.GetIndexPage(url)
			if err != nil {
				if errors.Is(err, helpers.ErrStatusNot200) {
					errorHandler <- "WRONG_STATUS_CODE"
					<-statusChannel
					errorOccurred = true
					return
				}
				errorHandler <- err
				<-statusChannel
				errorOccurred = true
				return
			}
			// now scan the page for possible folders
			folders := helpers.FilterDocumentForFolders(page)
			// now build the urls pointing into the folders for the
			// data types that will be checked next
			var urls []string
			for _, folder := range folders {
				url := fmt.Sprintf("%s/%s", url, folder)
				urls = append(urls, url)
			}
			dataTypeUrlChan <- urls
		}()
	}
	// wait for the wait group to finish
	waitGroup.Wait()
	// now check if any error occurred
	if errorOccurred {
		return
	}

	// now collect the urls
	var dataTypeUrls []string
	for len(dataTypeUrlChan) > 0 {
		urls := <-dataTypeUrlChan
		dataTypeUrls = append(dataTypeUrls, urls...)
	}

	// now check the index pages of the data types for their subdirectories
	waitGroup.Add(len(dataTypeUrls))
	dataUrlChan := make(chan []string, len(dataTypeUrls))
	for _, dataTypeUrl := range dataTypeUrls {
		dataTypeUrl := dataTypeUrl
		go func() {
			defer waitGroup.Done()
			log.Debug().Str("url", dataTypeUrl).Msg("requesting index page")
			page, err := helpers.GetIndexPage(dataTypeUrl)
			if err != nil {
				if errors.Is(err, helpers.ErrStatusNot200) {
					errorHandler <- "WRONG_STATUS_CODE"
					<-statusChannel
					errorOccurred = true
					return
				}
				errorHandler <- err
				<-statusChannel
				errorOccurred = true
				return
			}
			// now filter the page for folders and station files
			folders := helpers.FilterDocumentForFolders(page)
			// now build the array of possible folder and file urls
			var urls []string
			for _, folder := range folders {
				url := fmt.Sprintf("%s%s", dataTypeUrl, folder)
				urls = append(urls, url)
			}
			// now send the generated urls to the channel
			dataUrlChan <- urls
		}()
	}
	waitGroup.Wait()

	var dataUrls []string
	for len(dataUrlChan) > 0 {
		urls := <-dataUrlChan
		dataUrls = append(dataUrls, urls...)
	}
	waitGroup.Add(len(dataUrls))
	stationFileUrlChan := make(chan []string, len(dataUrls))
	for _, dataUrl := range dataUrls {
		dataUrl := dataUrl
		go func() {
			defer waitGroup.Done()
			log.Debug().Str("url", dataUrl).Msg("requesting index page")
			page, err := helpers.GetIndexPage(dataUrl)
			if err != nil {
				if errors.Is(err, helpers.ErrStatusNot200) {
					errorHandler <- "WRONG_STATUS_CODE"
					<-statusChannel
					errorOccurred = true
					return
				}
				errorHandler <- err
				<-statusChannel
				errorOccurred = true
				return
			}
			// now filter the page for folders and station files
			possibleFiles := helpers.FilterDocumentForFiles(page)
			// now build the array of possible folder and file urls
			var urls []string
			for _, file := range possibleFiles {
				if !strings.HasSuffix(file, "Beschreibung_Stationen.txt") {
					continue
				}
				url := fmt.Sprintf("%s%s", dataUrl, file)
				urls = append(urls, url)
			}
			// now send the generated urls to the channel
			stationFileUrlChan <- urls
		}()
	}
	waitGroup.Wait()

	var stationFileUrls []string
	for len(stationFileUrlChan) > 0 {
		urls := <-stationFileUrlChan
		stationFileUrls = append(stationFileUrls, urls...)
	}

	stationFileChan := make(chan []string, len(stationFileUrls))
	waitGroup.Add(len(stationFileUrls))
	for _, stationFileUrl := range stationFileUrls {
		stationFileUrl := stationFileUrl
		go func() {
			defer waitGroup.Done()
			log.Debug().Str("url", stationFileUrl).Msg("downloading station file")
			file, err := helpers.DownloadFile(stationFileUrl)
			if err != nil {
				if errors.Is(err, helpers.ErrStatusNot200) {
					errorHandler <- "WRONG_STATUS_CODE"
					<-statusChannel
					errorOccurred = true
					return
				}
				errorHandler <- err
				<-statusChannel
				errorOccurred = true
				return
			}
			stationFileChan <- []string{file.Name(), stationFileUrl}
		}()
	}
	waitGroup.Wait()

	// now iterate over the files to get the records and generate stations from
	// them
	stationMap := make(map[string]types.Station)
	for len(stationFileChan) > 0 {
		file := <-stationFileChan
		records, err := helpers.ParseStationFile(file[0])
		if err != nil {
			errorHandler <- err
			<-statusChannel
			return
		}
		for _, record := range records {
			stationID := record[0]
			fromDateString := record[1]
			untilDateString := record[2]
			fromDate, err := time.Parse(helpers.DateFormat_NoTime, fromDateString)
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
			untilDate, err := time.Parse(helpers.DateFormat_NoTime, untilDateString)
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
			station, stationAlreadyCreated := stationMap[stationID]
			if stationAlreadyCreated {
				resolution := types.Resolution(0)
				resolution.ParseStringWithSeparator(file[1], "/")

				dataType := types.DataType(0)
				dataType.ParseStringWithSeperator(file[1], "/")
				if dataType == 0 {
					panic(file[1])
				}
				// now construct the data capability
				c := types.Capability{
					DataType:       dataType,
					Resolution:     resolution,
					AvailableFrom:  fromDate,
					AvailableUntil: untilDate,
				}
				station.AddCapability(c)
				stationMap[stationID] = station

				continue
			}

			// since the station has not been created beforehand, get the needed
			// data from the records and create a station
			heightString := record[3]
			height, err := strconv.ParseFloat(heightString, 64)
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
			latString := record[4]
			lat, err := strconv.ParseFloat(latString, 64)
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
			longString := record[5]
			long, err := strconv.ParseFloat(longString, 64)
			if err != nil {
				errorHandler <- err
				<-statusChannel
				return
			}
			name := strings.Join(record[6:len(record)-2], " ")
			state := record[len(record)-2]

			location := geojson.NewPointGeometry([]float64{lat, long})

			s := types.Station{
				ID:       stationID,
				Name:     name,
				State:    state,
				Height:   height,
				Location: location,
			}

			resolution := types.Resolution(0)
			resolution.ParseStringWithSeparator(file[1], "/")

			dataType := types.DataType(0)
			dataType.ParseStringWithSeperator(file[1], "/")
			c := types.Capability{
				DataType:       dataType,
				Resolution:     resolution,
				AvailableFrom:  fromDate,
				AvailableUntil: untilDate,
			}
			s.AddCapability(c)
			stationMap[stationID] = s
		}
		os.Remove(file[0])
	}
	var stations []types.Station
	for _, station := range stationMap {
		err = station.UpdateHistoricalState()
		if err != nil {
			errorHandler <- err
			<-statusChannel
			return
		}
		stations = append(stations, station)
	}

	// encode the stations
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stations)

}

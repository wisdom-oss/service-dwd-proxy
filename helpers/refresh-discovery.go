package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
	geojson "github.com/paulmach/go.geojson"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"

	"github.com/wisdom-oss/service-dwd-proxy/globals"
	"github.com/wisdom-oss/service-dwd-proxy/types"
)

// BaseHost contains the opendata portals base host
const BaseHost = "https://opendata.dwd.de"

// BasePath contains the path to the climate observations
const BasePath = "/climate_environment/CDC/observations_germany/climate"

// client is the HTTP Client used to access the open data portal
var client = http.Client{}

// RunDiscovery queries the DWD Open Data server for the information about all
// stations
func RunDiscovery() {

	l := log.With().Str("part", "discovery").Logger()
	l.Info().Msg("checking status of discovery")
	discoveryContext := context.Background()
	defer globals.RedisClient.Set(discoveryContext, "dwd-discovery-running", "false", 0)
	discoveryRunning, err := globals.RedisClient.Get(discoveryContext, "dwd-discovery-running").Bool()
	if err != nil {
		if err == redis.Nil {
			globals.RedisClient.Set(discoveryContext, "dwd-discovery-running", "true", 0)
		} else {
			l.Error().Err(err).Msg("discovery of stations failed. service may return error")
			return
		}

	}

	if discoveryRunning {
		l.Warn().Msg("discovery already running. skipping this turn")
		return
	}
	l.Info().Msg("no running discovery detected. starting new one")
	// since there is no discovery running right now: execute one.
	// request the initial page of the index for the climate observations
	url := fmt.Sprintf("%s%s", BaseHost, BasePath)
	res, err := client.Get(url)
	if err != nil {
		l.Error().Err(err).Msg("discovery of stations failed. service may return error")
		return
	}
	if res.StatusCode != 200 {
		l.Error().Int("httpStatusCode", res.StatusCode).Msg("unexpected http code in response. discovery failed")
		return
	}
	// now parse the html tokens in on the page
	document, err := html.Parse(res.Body)
	if err != nil {
		l.Error().Err(err).Msg("discovery of stations failed. service may return error")
		return
	}
	// now resolve the folders that are listed in the document to be able to
	// find the currently available data resolutionFolders
	resolutionFolders := FilterDocumentForFolders(document)

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
			l.Debug().Str("url", url).Msg("requesting index page")
			page, err := GetIndexPage(url)
			if err != nil {
				l.Error().Err(err).Msg("discovery of stations failed. service may return error")
				errorOccurred = true
				return
			}
			// now scan the page for possible folders
			folders := FilterDocumentForFolders(page)
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
			l.Debug().Str("url", dataTypeUrl).Msg("requesting index page")
			page, err := GetIndexPage(dataTypeUrl)
			if err != nil {
				l.Error().Err(err).Msg("discovery of stations failed. service may return error")
				return
			}
			// now filter the page for folders and station files
			folders := FilterDocumentForFolders(page)
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
			l.Debug().Str("url", dataUrl).Msg("requesting index page")
			page, err := GetIndexPage(dataUrl)
			if err != nil {
				l.Error().Err(err).Msg("discovery of stations failed. service may return error")
				return
			}
			// now filter the page for folders and station files
			possibleFiles := FilterDocumentForFiles(page)
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
			l.Debug().Str("url", stationFileUrl).Msg("downloading station file")
			file, err := DownloadFile(stationFileUrl)
			if err != nil {
				l.Error().Err(err).Msg("discovery of stations failed. service may return error")
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
		records, err := ParseStationFile(file[0])
		if err != nil {
			l.Error().Err(err).Msg("discovery of stations failed. service may return error")
			return
		}
		for _, record := range records {
			stationID := record[0]
			fromDateString := record[1]
			untilDateString := record[2]
			fromDate, err := time.Parse(DateFormat_NoTime, fromDateString)
			if err != nil {
				l.Error().Err(err).Msg("discovery of stations failed. service may return error")
				return
			}
			untilDate, err := time.Parse(DateFormat_NoTime, untilDateString)
			if err != nil {
				l.Error().Err(err).Msg("unable to run discovery")
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
				l.Error().Err(err).Msg("unable to run discovery")
			}
			latString := record[4]
			lat, err := strconv.ParseFloat(latString, 64)
			if err != nil {
				l.Error().Err(err).Msg("unable to run discovery")
			}
			longString := record[5]
			long, err := strconv.ParseFloat(longString, 64)
			if err != nil {
				l.Error().Err(err).Msg("unable to run discovery")
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
			l.Error().Err(err).Msg("discovery of stations failed. service may return error")
			return
		}
		stations = append(stations, station)
	}

	// now write the stations and their data into the redis database
	stationBytes, err := json.Marshal(&stations)
	if err != nil {
		l.Error().Err(err).Msg("unable to marshal station list")
		return
	}
	var compressedByteBuffer bytes.Buffer
	brotliWriter := brotli.NewWriterLevel(&compressedByteBuffer, 9)
	writtenBytes, err := brotliWriter.Write(stationBytes)
	if err != nil {
		l.Error().Err(err).Msg("unable to compress json response. service may return error")
		return
	}
	err = brotliWriter.Flush()
	if err != nil {
		l.Error().Err(err).Msg("unable to flush brotli writer to buffer")
		return
	}
	fmt.Println("compressed response into bytes. size:", writtenBytes)
	err = globals.RedisClient.Set(discoveryContext, "dwd-station-list", compressedByteBuffer.Bytes(), 0).Err()
	if err != nil {
		l.Error().Err(err).Msg("unable to store station response. service may return error")
		return
	}
	l.Info().Msg("discovery data updated")
}

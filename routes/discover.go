package routes

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"

	"github.com/wisdom-oss/service-dwd-proxy/helpers"
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
	dataTypeUrls := make(chan []string, len(resolutionFolders))
	for _, resolutionFolder := range resolutionFolders {
		resolutionFolder := resolutionFolder
		go func() {
			// wait until this function exits to notify the wait group of this
			defer waitGroup.Done()
			// now build the url
			url := fmt.Sprintf("%s/%s", url, resolutionFolder)
			// now request the index page
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
			// now build the urls that will be checked next
			var urls []string
			for _, folder := range folders {
				url := fmt.Sprintf("%s%s", url, folder)
				urls = append(urls, url)
			}
			dataTypeUrls <- urls
		}()
	}
	// wait for the wait group to finish
	waitGroup.Wait()
	// now check if any error occurred
	if errorOccurred {
		return
	}

}

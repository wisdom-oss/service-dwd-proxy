package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"microservice/internal/dwd"
	types "microservice/types/v1"
	v1 "microservice/types/v1"
)

type timeseriesParameter struct {
	Start int `form:"from"`
	End   int `form:"until"`
}

func Timeseries(c *gin.Context) {
	var params timeseriesParameter
	err := c.BindQuery(&params)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	datapoint := types.DataType(0)
	datapoint.ParseString(c.Param("dataType"))

	if datapoint == 0 {
		errUnsupportedDatapoint.Emit(c)
		return
	}

	resolution := types.Resolution(0)
	resolution.ParseString(c.Param("resolution"))

	url := fmt.Sprintf("%s/%s/%s/%s", DWD_OpenData_Host, DWD_OpenData_Base, resolution, datapoint)

	page, err := dwd.LoadIndexPage(url)
	if err != nil {
		c.Abort()

		if errors.Is(err, dwd.ErrNotFound) {
			errUnsupportedDatapoint.Emit(c)
			return
		}

		_ = c.Error(fmt.Errorf("unable to load index page: %w", err))
		return
	}

	fileUrls := []string{}

	folderUrls := dwd.GetFolderURLs(page, url)
	for _, folderUrl := range folderUrls {
		page, err := dwd.LoadIndexPage(folderUrl)
		if err != nil {
			c.Abort()

			if errors.Is(err, dwd.ErrNotFound) {
				errUnsupportedDatapoint.Emit(c)
				return
			}

			_ = c.Error(fmt.Errorf("unable to load index page: %w", err))
			return
		}

		files := dwd.FilterDocumentForFiles(page)
		for _, file := range files {
			if strings.Contains(file, c.Param("stationID")) {
				fileUrls = append(fileUrls, fmt.Sprintf("%s%s", folderUrl, file))
			}
		}

		if len(fileUrls) == 0 {
			c.Abort()
			c.Status(999)
			return
		}
	}

	metadataFiles := []string{}
	dataFiles := []string{}

	for _, url := range fileUrls {
		file, err := dwd.Download(url)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			return
		}

		if strings.Contains(url, "meta_data") {
			metadataFiles = append(metadataFiles, file)
		} else {
			dataFiles = append(dataFiles, file)
		}

	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var parameters []v1.TimeseriesField
	var data []map[string]any

	wg.Add(2)
	go func(fileNames []string) {
		defer wg.Done()

		for _, name := range fileNames {
			metdataFields, _ := dwd.ParseMetadataArchive(name)
			mutex.Lock()
			parameters = append(parameters, metdataFields...)
			mutex.Unlock()
		}

	}(metadataFiles)
	go func(fileNames []string) {
		defer wg.Done()

		for _, name := range fileNames {
			dps, metadata, _ := dwd.ParseDataFile(name, [2]time.Time{time.Unix(int64(params.Start), 0), time.Unix(int64(params.End), 0)}) //nolint:lll
			data = append(data, dps...)
			mutex.Lock()
			parameters = append(parameters, metadata...)
			mutex.Unlock()
		}

	}(dataFiles)
	wg.Wait()
	type res struct {
		Data []map[string]interface{} `json:"timeseries"`
		Meta []types.TimeseriesField  `json:"metadata"`
	}
	response := res{
		Data: data,
		Meta: parameters,
	}

	c.JSON(http.StatusOK, response)
}

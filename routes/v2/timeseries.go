package v2

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/wisdom-oss/common-go/v3/types"

	dwd "microservice/internal/dwd/v2"
	v2 "microservice/types/v2"
)

var errUnknownDatabase = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.5.5",
	Status: http.StatusNotFound,
	Title:  "Database Unknown",
	Detail: "The requested database is not known to this proxy",
}

var errDatabaseUnreachable = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.6.4",
	Status: http.StatusServiceUnavailable,
	Title:  "Database Unreachable",
	Detail: "The database is currently not able to handle requests",
}

var errUnknownProduct = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Unknown Product",
	Detail: "The supplied product is not available in this database",
}

var errUnknownGranularity = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Unknown Granularity",
	Detail: "The supplied granularity is unknown",
}

var errUnsupportedGranularity = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Unsupported Product/Granularity",
	Detail: "The supplied product/granularity combination is not supported",
}

var errStationValidationFailed = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.6.1",
	Status: http.StatusInternalServerError,
	Title:  "Station Validation Failed",
	Detail: "Unable to determine if the supplied station is available for this product/granularity combination",
}

var errStationNotAvailable = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Station Not Available",
	Detail: "The supplied station is not available for the selected product/granularity combination",
}

var errTimeseriesStartTooEarly = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Timeseries Starts Too Early",
	Detail: "The selected start of the timeseries is before the station has delivered data",
}

var errTimeseriesEndTooLate = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Timeseries End Too Late",
	Detail: "The selected end of the timeseries is after the station has been decommissioned",
}

var errTimeseriesParseError = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Timeseries Parsing Error",
	Detail: "The provided timestamps could not be parsed from the query parameters",
}

var errTimeseriesBoundaryError = types.ServiceError{
	Type:   "https://datatracker.ietf.org/doc/html/rfc9110#section-15.5.1",
	Status: http.StatusBadRequest,
	Title:  "Timeseries Boundaries Switched",
	Detail: "The boundaries of the timeseries are not valid (start is after end)",
}

func Timeseries(c *gin.Context) {
	database := c.Param("database")
	databaseKeys := make([]string, 0, len(dwd.Databases))
	for k := range dwd.Databases {
		databaseKeys = append(databaseKeys, k)
	}

	if !slices.Contains(databaseKeys, database) {
		c.Abort()
		errUnknownDatabase.Emit(c)
		return
	}

	res, err := http.Get(dwd.Databases[database])
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	if res.StatusCode != http.StatusOK {
		c.Abort()
		errDatabaseUnreachable.Emit(c)
		return
	}

	// check if the product is supported
	p := c.Param("product")
	product := dwd.Product(0)
	if err := product.Parse(p); err != nil {
		c.Abort()
		errUnknownProduct.Emit(c)
		return
	}

	// check if the granularity is supported for the product
	g := c.Param("granularity")
	granularity := dwd.Granularity(0)
	if err := granularity.Parse(g); err != nil {
		c.Abort()
		errUnknownGranularity.Emit(c)
		return
	}

	if !slices.Contains(dwd.AvailableClimateObservationProducts[granularity], product) {
		c.Abort()
		errUnsupportedGranularity.Emit(c)
		return
	}

	// now request the station list for the product
	stations, err := dwd.DiscoverStations(dwd.Databases[database], granularity, product)
	if err != nil {
		c.Abort()
		errStationValidationFailed.Emit(c)
		return
	}

	var station v2.Station
	if !slices.ContainsFunc(stations, func(s v2.Station) bool {
		if s.ID == c.Param("stationID") {
			station = s
			return true
		}
		return false
	}) {
		c.Abort()
		errStationNotAvailable.Emit(c)
		return
	}

	var requestedRange struct {
		Start time.Time `form:"start"`
		End   time.Time `form:"end"`
	}
	if err := c.ShouldBindQuery(&requestedRange); err != nil {
		c.Abort()
		errTimeseriesParseError.Emit(c)
		return
	}

	dataAvailableFrom := station.SupportedProducts[product][granularity]

	if requestedRange.Start.IsZero() && requestedRange.End.IsZero() {
		goto startDownload
	}

	if requestedRange.Start.After(requestedRange.End) && !requestedRange.End.IsZero() {
		c.Abort()
		errTimeseriesBoundaryError.Emit(c)
		return
	}

	if requestedRange.Start.Before(dataAvailableFrom.Start) {
		c.Abort()
		errTimeseriesStartTooEarly.Emit(c)
		return
	}

	if requestedRange.End.After(dataAvailableFrom.End) {
		c.Abort()
		errTimeseriesEndTooLate.Emit(c)
		return
	}

startDownload:

	dataFiles, descriptionFiles, err := dwd.DownloadFiles(database, station.ID, product, granularity)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var series v2.Timeseries

	for _, descriptionFile := range descriptionFiles {
		var file v2.File

		if strings.HasPrefix(descriptionFile[0], "BESCHREIBUNG") {
			file.Name = "[DE] Datensatzbeschreibung"
		}

		if strings.HasPrefix(descriptionFile[0], "DESCRIPTION") {
			file.Name = "[EN] Dataset Description"
		}

		if file.Name == "" {
			file.Name = strings.Trim(strings.SplitAfterN(descriptionFile[0], ".", 2)[0], ".") //nolint:mnd
		}

		f, err := os.Open(descriptionFile[1])
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			return
		}

		mime, err := mimetype.DetectReader(f)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			return
		}
		file.MimeType = mime.String()

		var buf bytes.Buffer
		enc := base64.NewEncoder(base64.StdEncoding, &buf)
		_, _ = f.Seek(0, io.SeekStart)
		_, err = io.Copy(enc, f)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			return
		}
		_ = enc.Close()

		file.Content = buf.String()

		series.DescriptionFiles = append(series.DescriptionFiles, file)

	}

	allDatapoints := make([]v2.Datapoint, 0)
	allMetadata := make([]v2.FieldMetadata, 0)

	for _, dataFile := range dataFiles {
		datapoints, metadata, err := dwd.HandleArchive(dataFile)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			return
		}
		allDatapoints = append(allDatapoints, datapoints...)
		allMetadata = append(allMetadata, metadata...)
	}

	for {
		if slices.IsSortedFunc(allDatapoints, func(this, other v2.Datapoint) int {
			return this.Timestamp.Compare(other.Timestamp)
		}) {
			break
		}

		slices.SortFunc(allDatapoints, func(this, other v2.Datapoint) int {
			return this.Timestamp.Compare(other.Timestamp)
		})
	}

	series.Datapoints = allDatapoints
	series.Metadata = allMetadata

	c.JSON(http.StatusOK, series)

}

package parser

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	_ "time/tzdata"

	"github.com/twpayne/go-geom"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	v2 "microservice/types/v2"
)

// columnCount is the number of columns a correctly formatted station list.
const columnCount = 9
const coordinateSRID = 4326

const (
	idx_StationID = iota
	idx_DataStartDate
	idx_DataEndDate
	idx_StationHeight
	idx_Latitude
	idx_Longitude
	idx_StationName
	idx_State
	idx_Fee
)

const (
	df_DayOnly  = "20060102"
	df_HourOnly = "2006010215"
	df_Full     = "200601021504"
)

func ParseStationList(r io.Reader) (stations []v2.Station, dates [][2]time.Time, err error) {
	csvReader := csv.NewReader(transform.NewReader(r, charmap.Windows1252.NewDecoder()))
	csvReader.Comma = ' '
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	for num, data := range rows {
		// skip the headers in the file
		if num <= 1 {
			continue
		}

		var cleanedData []string

		if len(data) > columnCount+1 { // plus one to account for the empty delimiter
			cleanedData = append(cleanedData, data[:6]...)
			additionalNameParts := len(data) - columnCount
			nameParts := data[6 : 6+additionalNameParts]
			name := strings.Join(nameParts, " ")
			cleanedData = append(cleanedData, name)
			cleanedData = append(cleanedData, data[6+additionalNameParts:]...)
		} else {
			cleanedData = data
		}

		var longitude, latitude, height float64
		longitude, err = strconv.ParseFloat(cleanedData[idx_Longitude], 64)
		if err != nil {
			return nil, nil, err
		}

		latitude, err = strconv.ParseFloat(cleanedData[idx_Latitude], 64)
		if err != nil {
			return nil, nil, err
		}

		height, err = strconv.ParseFloat(cleanedData[idx_StationHeight], 64)
		if err != nil {
			return nil, nil, err
		}

		location := geom.NewPointFlat(geom.XYZ, []float64{longitude, latitude, height})
		location.SetSRID(coordinateSRID)

		station := v2.Station{
			ID:       cleanedData[idx_StationID],
			Name:     cleanedData[idx_StationName],
			Location: location,
		}

		var startDate, endDate time.Time
		switch len(cleanedData[idx_DataStartDate]) {
		case len(df_Full):
			startDate, err = time.Parse(df_Full, cleanedData[idx_DataStartDate])
		case len(df_HourOnly):
			startDate, err = time.Parse(df_HourOnly, cleanedData[idx_DataStartDate])
		case len(df_DayOnly):
			startDate, err = time.Parse(df_DayOnly, cleanedData[idx_DataStartDate])
		default:
			return nil, nil, errors.New("unsupported date string for data start")
		}
		if err != nil {
			return nil, nil, err
		}

		switch len(strings.TrimSpace(cleanedData[idx_DataEndDate])) {
		case len(df_Full):
			endDate, err = time.Parse(df_Full, (cleanedData[idx_DataEndDate]))
		case len(df_HourOnly):
			endDate, err = time.Parse(df_HourOnly, (cleanedData[idx_DataEndDate]))
		case len(df_DayOnly):
			endDate, err = time.Parse(df_DayOnly, (cleanedData[idx_DataEndDate]))
		default:
			return nil, nil, errors.New("unsupported date string for data end")

		}
		if err != nil {
			return nil, nil, err
		}

		mez, err := time.LoadLocation("Etc/GMT-1")
		if err != nil {
			return nil, nil, err
		}

		if startDate.Year() < 2000 { //nolint:mnd
			startDate = startDate.In(mez)
		}

		if endDate.Year() < 2000 { //nolint:mnd

			endDate = startDate.In(mez)
		}

		stations = append(stations, station)
		dates = append(dates, [2]time.Time{startDate, endDate})

	}
	return
}

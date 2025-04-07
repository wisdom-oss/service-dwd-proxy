package parser

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"

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

func ParseStationList(r io.Reader) (stations []v2.Station, err error) {
	csvReader := csv.NewReader(transform.NewReader(r, charmap.Windows1252.NewDecoder()))
	csvReader.Comma = ' '
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
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
			return nil, err
		}

		latitude, err = strconv.ParseFloat(cleanedData[idx_Latitude], 64)
		if err != nil {
			return nil, err
		}

		height, err = strconv.ParseFloat(cleanedData[idx_StationHeight], 64)

		location := geom.NewPointFlat(geom.XYZ, []float64{longitude, latitude, height})
		location.SetSRID(coordinateSRID)

		station := v2.Station{
			ID:       cleanedData[idx_StationID],
			Name:     cleanedData[idx_StationName],
			Location: location,
		}

		stations = append(stations, station)

	}
	return
}

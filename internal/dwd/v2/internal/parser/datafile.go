package parser

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"microservice/internal/dwd/v2/dwdTypes"
	v2 "microservice/types/v2"
)

const (
	filePrefix_Parameters    = "Metadaten_Parameter_"
	filePrefix_MissingValues = "Metadaten_Fehlwerte_"
	filePrefix_DataFile      = "produkt_"
)

const (
	metadataFieldName_Parameter   = "Parameter"
	metadataFieldName_Description = "Parameterbeschreibung"
	metadataFieldName_Unit        = "Einheit"
	metadataFieldName_From        = "Von_Datum"
	metadataFieldName_Until       = "Bis_Datum"
)

const (
	dataFieldName_StationID      = "STATIONS_ID"
	dataFieldName_Date           = "MESS_DATUM"
	dataFieldPrefix_QualityLevel = "QN"
	dataField_EndOfRow           = "eor"
)

const (
	fieldName_MissingValues_StationID  = "Stations_ID"
	fieldName_MissingValues_Start      = "Von_Datum"
	fieldName_MissingValues_End        = "Bis_Datum"
	fieldName_MissingValues_Count      = "Anzahl_Fehlwerte"
	fieldName_MissingValues_Parameters = "Parameter"

	df_missing_day_time = "02.01.2006-15:04"
)

func ReadArchive(path string) (datapoints []v2.Datapoint, metadata []v2.FieldMetadata, err error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return nil, nil, err
	}

	var parsedDatapoints []v2.Datapoint
	var generatedDatapoints []v2.Datapoint

	for _, file := range archive.File {
		if !strings.HasSuffix(file.Name, ".txt") {
			continue
		}

		if strings.HasPrefix(file.Name, filePrefix_Parameters) {
			metadata, err = parseMetadataFile(file)
			if err != nil {
				return nil, nil, err
			}
			continue
		}

		if strings.HasPrefix(file.Name, filePrefix_MissingValues) {
			generatedDatapoints, err = generateMissingDatapoints(file)
			if err != nil {
				return nil, nil, err
			}
			datapoints = append(datapoints, generatedDatapoints...)
			continue
		}

		if !strings.HasPrefix(file.Name, filePrefix_DataFile) {
			continue
		}

		parsedDatapoints, err = parseDatapointFile(file)
		if err != nil {
			return nil, nil, err
		}

	}

	if metadata == nil || datapoints == nil {
		return
	}

	units := make(map[string]string)
	for _, metadataField := range metadata {
		units[metadataField.Name] = metadataField.Unit
	}

	for _, dp := range parsedDatapoints {
		l := units[dp.Label]
		dp.Unit = &l
		datapoints = append(datapoints, dp)
	}

	return
}

func parseMetadataFile(compressedFile *zip.File) (metadata []v2.FieldMetadata, err error) {
	f, err := compressedFile.Open()
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(transform.NewReader(f, charmap.Windows1252.NewDecoder().Transformer))
	reader.TrimLeadingSpace = true
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := lines[0]

	for idx, line := range lines {
		if idx == 0 || len(line) != len(header) { // exclude the header and all lines not matching the header length
			continue
		}

		m := v2.FieldMetadata{
			Name:        line[slices.Index(header, metadataFieldName_Parameter)],
			Description: line[slices.Index(header, metadataFieldName_Description)],
			Unit:        line[slices.Index(header, metadataFieldName_Unit)],
		}

		var from time.Time
		fromDateString := line[slices.Index(header, metadataFieldName_From)]
		switch len(fromDateString) {
		case len(df_Full):
			from, err = time.Parse(df_Full, fromDateString)
		case len(df_HourOnly):
			from, err = time.Parse(df_HourOnly, fromDateString)
		case len(df_DayOnly):
			from, err = time.Parse(df_DayOnly, fromDateString)
		default:
			return nil, errors.New("unsupported datetime format")
		}
		if err != nil {
			return nil, err
		}

		m.ValidFrom = from

		var until time.Time
		untilDateString := line[slices.Index(header, metadataFieldName_Until)]
		switch len(fromDateString) {
		case len(df_Full):
			until, err = time.Parse(df_Full, untilDateString)
		case len(df_HourOnly):
			until, err = time.Parse(df_HourOnly, untilDateString)
		case len(df_DayOnly):
			until, err = time.Parse(df_DayOnly, untilDateString)
		default:
			return nil, errors.New("unsupported datetime format")
		}
		if err != nil {
			return nil, err
		}

		m.ValidUntil = until

		metadata = append(metadata, m)

	}

	return metadata, nil

}

func parseDatapointFile(compressedFile *zip.File) (datapoints []v2.Datapoint, err error) {
	f, err := compressedFile.Open()
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(transform.NewReader(f, charmap.Windows1252.NewDecoder().Transformer))
	reader.TrimLeadingSpace = true
	reader.Comma = ';'

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := lines[0]
	ql_idx := -1
	datacolidxs := make(map[int]string)

	for idx, rowHead := range header {
		if rowHead == dataFieldName_Date || rowHead == dataFieldName_StationID || rowHead == dataField_EndOfRow {
			continue
		}

		if strings.HasPrefix(rowHead, dataFieldPrefix_QualityLevel) {
			ql_idx = idx
			continue
		}

		datacolidxs[idx] = rowHead
	}

	for idx, line := range lines {
		if idx == 0 {
			continue
		}

		var date time.Time
		dateString := line[slices.Index(header, dataFieldName_Date)]
		switch len(dateString) {
		case len(df_Full):
			date, err = time.Parse(df_Full, dateString)
		case len(df_HourOnly):
			date, err = time.Parse(df_HourOnly, dateString)
		case len(df_DayOnly):
			date, err = time.Parse(df_DayOnly, dateString)
		default:
			return nil, errors.New("unsupported datetime format")
		}
		if err != nil {
			return nil, err
		}

		if date.Year() < 2000 { //nolint:mnd
			mez, err := time.LoadLocation("Etc/GMT-1")
			if err != nil {
				return nil, err
			}
			date = date.In(mez)
		}

		for idx, name := range datacolidxs {
			p := v2.Datapoint{
				Label:     name,
				Timestamp: date,
			}

			val := line[idx]
			floatValue, err := strconv.ParseFloat(val, 64)
			if err != nil {
				p.Value = val
			} else {
				p.Value = floatValue
			}

			qualityLevelStr := line[ql_idx]
			qualityLevel := dwdTypes.QualityFlag(0)
			qlInt, err := strconv.ParseInt(qualityLevelStr, 10, 64)
			if err != nil {
				return nil, err
			}

			err = qualityLevel.Parse(qlInt)
			if err != nil {
				return nil, err
			}

			p.QualityLevel = &qualityLevel

			datapoints = append(datapoints, p)
		}
	}

	return
}

func generateMissingDatapoints(compressedFile *zip.File) (generatedDatapoints []v2.Datapoint, err error) {
	f, err := compressedFile.Open()
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(transform.NewReader(f, charmap.Windows1252.NewDecoder().Transformer))
	reader.TrimLeadingSpace = true
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := lines[0]
	fmt.Println(header)

	for idx, line := range lines {
		if idx == 0 || line[0] == fieldName_MissingValues_StationID || idx == len(lines)-1 {
			continue
		}

		missingDataStart, err := time.Parse(df_missing_day_time, line[slices.Index(header, fieldName_MissingValues_Start)])
		if err != nil {
			return nil, err
		}

		missingDataEnd, err := time.Parse(df_missing_day_time, line[slices.Index(header, fieldName_MissingValues_End)])
		if err != nil {
			return nil, err
		}

		missingDatapointCount, err := strconv.ParseInt(line[slices.Index(header, fieldName_MissingValues_Count)], 10, 64)
		if err != nil {
			return nil, err
		}

		if missingDatapointCount == 1 {
			generatedDatapoints = append(generatedDatapoints, v2.Datapoint{
				Label:        line[slices.Index(header, fieldName_MissingValues_Parameters)],
				Value:        nil,
				Timestamp:    missingDataStart,
				Unit:         nil,
				QualityLevel: nil,
			})
			continue
		}

		timeDiff := missingDataEnd.Sub(missingDataStart)
		timeStep := timeDiff.Nanoseconds() / (missingDatapointCount - 1)
		for i := range missingDatapointCount - 1 {
			offset := time.Duration(i * timeStep)
			generatedDatapoints = append(generatedDatapoints, v2.Datapoint{
				Label:        line[slices.Index(header, fieldName_MissingValues_Parameters)],
				Value:        nil,
				Timestamp:    missingDataStart.Add(offset),
				Unit:         nil,
				QualityLevel: nil,
			})
		}

	}

	return generatedDatapoints, nil
}

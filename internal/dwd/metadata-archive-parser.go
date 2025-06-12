package dwd

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"slices"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	v1 "microservice/types/v1"
)

func ParseMetadataArchive(path string) (parameters []v1.TimeseriesField, err error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	for _, file := range zipFile.File {
		if !strings.HasSuffix(file.Name, "txt") || !strings.Contains(file.Name, "Metadaten_Parameter") {
			fmt.Println("skipping file", file.Name)
			continue
		}
		fmt.Println("reading file", file.Name)
		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		windows1252Reader := transform.NewReader(f, charmap.Windows1252.NewDecoder())
		csvReader := csv.NewReader(windows1252Reader)
		csvReader.TrimLeadingSpace = true
		csvReader.Comma = ';'
		csvReader.FieldsPerRecord = -1

		lines, err := csvReader.ReadAll()
		if err != nil {
			return nil, err
		}
		// split the lines into the header and the contents
		header := lines[0]
		content := lines[1 : len(lines)-2]

		parameters, err = parseMetadataFileContents(header, content)
		if err != nil {
			return nil, err
		}
	}

	return parameters, nil
}

func parseMetadataFileContents(header []string, contentLines [][]string) (parameters []v1.TimeseriesField, err error) {
	for _, contentLine := range contentLines {
		var timeseriesField v1.TimeseriesField
		timeseriesField.Name = contentLine[slices.Index(header, ParameterNameKey)]
		timeseriesField.Description = contentLine[slices.Index(header, ParameterDescriptionKey)]
		timeseriesField.Unit = contentLine[slices.Index(header, ParameterUnitKey)]

		// now try to parse the date range in which the field is available
		fromDateString := contentLine[slices.Index(header, ParameterFromDateKey)]
		fromDate, err := time.Parse(DateFormat_DateTimeFull, fromDateString)
		if err != nil {
			fromDate, err = time.Parse(DateFormat_DateTimeHourOnly, fromDateString)
			if err != nil {
				fromDate, err = time.Parse(DateFormat_NoTime, fromDateString)
				if err != nil {
					return nil, err
				}
			}
		}

		untilDateString := contentLine[slices.Index(header, ParameterUntilDateKey)]
		untilDate, err := time.Parse(DateFormat_DateTimeFull, untilDateString)
		if err != nil {
			untilDate, err = time.Parse(DateFormat_DateTimeHourOnly, untilDateString)
			if err != nil {
				untilDate, err = time.Parse(DateFormat_NoTime, untilDateString)
				if err != nil {
					return nil, err
				}
			}
		}

		timeseriesField.AvailableFrom = fromDate
		timeseriesField.AvailableUntil = untilDate

		parameters = append(parameters, timeseriesField)
	}
	return parameters, nil
}

package dwd

import (
	"archive/zip"
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	v1 "microservice/types/v1"
)

func ParseDataFile(path string, timeRange [2]time.Time) (datasets []map[string]interface{}, parameters []v1.TimeseriesField, err error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, nil, err
	}
	for _, file := range zipFile.File {
		if strings.HasSuffix(file.Name, ".txt") && strings.HasPrefix(file.Name, "Metadaten_Parameter") {
			f, err := file.Open()
			if err != nil {
				return nil, nil, err
			}
			windows1252Reader := transform.NewReader(f, charmap.Windows1252.NewDecoder())
			csvReader := csv.NewReader(windows1252Reader)
			csvReader.TrimLeadingSpace = true
			csvReader.Comma = ';'
			csvReader.FieldsPerRecord = -1

			lines, err := csvReader.ReadAll()
			if err != nil {
				return nil, nil, err
			}
			// split the lines into the header and the contents
			header := lines[0]
			content := lines[1 : len(lines)-2]

			parameters, err = parseMetadataFileContents(header, content)
			if err != nil {
				return nil, nil, err
			}
		}
		if !strings.HasSuffix(file.Name, ".txt") || !strings.Contains(file.Name, "produkt") {
			continue
		}
		f, err := file.Open()
		if err != nil {
			return nil, nil, err
		}
		windows1252Reader := transform.NewReader(f, charmap.Windows1252.NewDecoder())
		csvReader := csv.NewReader(windows1252Reader)
		csvReader.TrimLeadingSpace = true
		csvReader.Comma = ';'
		csvReader.FieldsPerRecord = -1

		lines, err := csvReader.ReadAll()
		if err != nil {
			return nil, nil, err
		}
		header := lines[0]
		// remove the station id from the header since it is not needed
		header = header[1 : len(header)-1]
		// now get the start date from which the datasets shall be included
		dataStartDate := timeRange[0]
		dataEndDate := timeRange[1]
		for _, data := range lines[2:] {
			// same here
			data = data[1 : len(data)-1]
			dataset := make(map[string]interface{})
			// now parse the measurement date
			measurementTime, err := time.Parse(DateFormat_DateTimeFull, data[0])
			if err != nil {
				measurementTime, err = time.Parse(DateFormat_DateTimeHourOnly, data[0])
				if err != nil {
					measurementTime, err = time.Parse(DateFormat_NoTime, data[0])
					if err != nil {
						return nil, nil, err
					}
				}
			}
			if !dataStartDate.IsZero() && measurementTime.Before(dataStartDate) {
				continue
			}
			if !dataEndDate.IsZero() && measurementTime.After(dataEndDate) {
				continue
			}
			dataset["ts"] = measurementTime
			for i := 1; i < len(data); i++ {
				// now try to parse every value into a float for easier
				// usage in other applications
				floatValue, err := strconv.ParseFloat(data[i], 64)
				if err != nil {
					dataset[header[i]] = data[i]
				} else {
					dataset[header[i]] = floatValue
				}
			}
			datasets = append(datasets, dataset)
		}
	}

	return datasets, parameters, nil
}

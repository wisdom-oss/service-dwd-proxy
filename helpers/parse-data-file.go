package helpers

import (
	"archive/zip"
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const measurementDateTimeFormat = "200601021504"

func ParseDataFile(path string) (datasets []map[string]interface{}, err error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	for _, file := range zipFile.File {
		if !strings.HasSuffix(file.Name, ".txt") && !strings.Contains(file.Name, "produkt") {
			continue
		}
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
		header := lines[0]
		// remove the station id from the header since it is not needed
		header = header[1 : len(header)-1]
		for _, data := range lines[2:] {
			// same here
			data = data[1 : len(data)-1]
			dataset := make(map[string]interface{})
			// now parse the measurement date
			measurementTime, err := time.Parse(measurementDateTimeFormat, data[0])
			if err != nil {
				return nil, err
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

	return datasets, nil
}

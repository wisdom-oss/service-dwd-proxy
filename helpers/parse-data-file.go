package helpers

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func ParseDataFile(path string) ([]map[string]interface{}, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	for _, file := range zipFile.File {
		if !strings.HasSuffix(file.Name, ".txt") {
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
		header = header[2 : len(header)-1]
		fmt.Println(header)
	}

	return nil, nil
}

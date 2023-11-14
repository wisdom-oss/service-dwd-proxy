package helpers

import (
	"encoding/csv"
	"os"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// ParseStationFile returns an array containing the single entries of a station
// file
func ParseStationFile(path string) (lines [][]string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	windows1252Reader := transform.NewReader(f, charmap.Windows1252.NewDecoder())

	csvReader := csv.NewReader(windows1252Reader)
	csvReader.Comma = ' '
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = -1
	lines, err = csvReader.ReadAll()
	return lines[2:], err
}

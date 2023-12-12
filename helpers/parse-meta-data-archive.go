package helpers

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/wisdom-oss/service-dwd-proxy/types"
)

// ParseMetadataArchive reads through the archive containing the parameter
// metadata and returns it.
func ParseMetadataArchive(path string) (parameters []types.TimeseriesField, err error) {
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

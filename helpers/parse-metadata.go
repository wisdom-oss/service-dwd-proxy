package helpers

import (
	"slices"
	"time"

	"github.com/wisdom-oss/service-dwd-proxy/types"
)

func parseMetadataFileContents(header []string, contentLines [][]string) (parameters []types.TimeseriesField, err error) {
	for _, contentLine := range contentLines {
		var timeseriesField types.TimeseriesField
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

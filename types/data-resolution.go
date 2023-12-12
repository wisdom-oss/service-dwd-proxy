package types

import (
	"encoding/json"
	"net/url"
	"strings"
)

type Resolution uint16

const (
	OneMinute = 1 << iota
	FiveMinutes
	TenMinutes
	Hourly
	SubDaily
	Daily
	Monthly
	Annually
	MultiAnnually
)

// String returns the value needed to
func (r Resolution) String() string {
	var result string
	if r&OneMinute != 0 {
		result += "1_minute,"
	}

	if r&FiveMinutes != 0 {
		result += "5_minutes,"
	}

	if r&TenMinutes != 0 {
		result += "10_minutes,"
	}

	if r&Hourly != 0 {
		result += "hourly,"
	}

	if r&SubDaily != 0 {
		result += "subdaily,"
	}

	if r&Daily != 0 {
		result += "daily,"
	}

	if r&Monthly != 0 {
		result += "monthly,"
	}

	if r&Annually != 0 {
		result += "annual,"
	}

	if r&MultiAnnually != 0 {
		result += "multi_annual,"
	}

	return strings.Trim(result, ",")
}

// ParseString takes an input string and tries to convert the string into a
// resolution.
// If multiple values are supposed to be in a string, the function expects
// the values to be separated by a comma.
// If another separator shall be used, please use the ParseStringWithSeparator
// function
func (r *Resolution) ParseString(s string) {
	var parsedResolution Resolution
	parts := strings.Split(s, ",")
	for _, part := range parts {
		switch part {
		case "1_minute":
			parsedResolution = parsedResolution | OneMinute
			break
		case "5_minutes":
			parsedResolution = parsedResolution | FiveMinutes
			break
		case "10_minutes":
			parsedResolution = parsedResolution | TenMinutes
			break
		case "hourly":
			parsedResolution = parsedResolution | Hourly
			break
		case "subdaily":
			parsedResolution = parsedResolution | SubDaily
			break
		case "daily":
			parsedResolution = parsedResolution | Daily
			break
		case "monthly":
			parsedResolution = parsedResolution | Monthly
			break
		case "annual":
			parsedResolution = parsedResolution | Annually
			break
		case "multi_annual":
			parsedResolution = parsedResolution | MultiAnnually
			break
		}
	}
	*r = parsedResolution
}

// ParseStringWithSeparator takes an input string and a separator and
// tries to convert the string into a resolution.
// If multiple values are supposed to be in a string, the function expects
// the values to be separated by a comma.
// If a comma is used as separator, please use the ParseString
// function as it expects this by default
func (r *Resolution) ParseStringWithSeparator(s, sep string) {
	var parsedResolution Resolution
	parts := strings.Split(s, sep)
	for _, part := range parts {
		switch part {
		case "1_minute":
			parsedResolution = parsedResolution | OneMinute
			break
		case "5_minutes":
			parsedResolution = parsedResolution | FiveMinutes
			break
		case "10_minutes":
			parsedResolution = parsedResolution | TenMinutes
			break
		case "hourly":
			parsedResolution = parsedResolution | Hourly
			break
		case "subdaily":
			parsedResolution = parsedResolution | SubDaily
			break
		case "daily":
			parsedResolution = parsedResolution | Daily
			break
		case "monthly":
			parsedResolution = parsedResolution | Monthly
			break
		case "annual":
			parsedResolution = parsedResolution | Annually
			break
		case "multi_annual":
			parsedResolution = parsedResolution | MultiAnnually
			break
		}
	}
	*r = parsedResolution
}

// ParseUrlValues takes the query parameters from the url object and the
// key under which the resolutions are available
func (r *Resolution) ParseUrlValues(q url.Values, key string) {
	var parsedResolution Resolution
	parts := q[key]
	for _, part := range parts {
		switch part {
		case "1_minute":
			parsedResolution = parsedResolution | OneMinute
			break
		case "5_minutes":
			parsedResolution = parsedResolution | FiveMinutes
			break
		case "10_minutes":
			parsedResolution = parsedResolution | TenMinutes
			break
		case "hourly":
			parsedResolution = parsedResolution | Hourly
			break
		case "subdaily":
			parsedResolution = parsedResolution | SubDaily
			break
		case "daily":
			parsedResolution = parsedResolution | Daily
			break
		case "monthly":
			parsedResolution = parsedResolution | Monthly
			break
		case "annual":
			parsedResolution = parsedResolution | Annually
			break
		case "multi_annual":
			parsedResolution = parsedResolution | MultiAnnually
			break
		}
	}
	*r = parsedResolution
}

func (r Resolution) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Resolution) UnmarshalJSON(src []byte) error {
	var resolutions []string
	if err := json.Unmarshal(src, &resolutions); err != nil {
		return err
	}
	for _, resolution := range resolutions {
		r.ParseString(resolution)
	}
	return nil
}

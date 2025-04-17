package v1

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

// String returns the value needed to.
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

// nolint: goconst
func (r *Resolution) ParseString(s string) {
	var parsedResolution Resolution
	parts := strings.Split(s, ",")
	for _, part := range parts {
		switch part {
		case "1_minute":
			parsedResolution |= OneMinute
		case "5_minutes":
			parsedResolution |= FiveMinutes
		case "10_minutes":
			parsedResolution |= TenMinutes
		case "hourly":
			parsedResolution |= Hourly
		case "subdaily":
			parsedResolution |= SubDaily
		case "daily":
			parsedResolution |= Daily
		case "monthly":
			parsedResolution |= Monthly
		case "annual":
			parsedResolution |= Annually
		case "multi_annual":
			parsedResolution |= MultiAnnually
		}
	}
	*r = parsedResolution
}

// nolint: goconst
func (r *Resolution) ParseStringWithSeparator(s, sep string) {
	var parsedResolution Resolution
	parts := strings.Split(s, sep)
	for _, part := range parts {
		switch part {
		case "1_minute":
			parsedResolution |= OneMinute

		case "5_minutes":
			parsedResolution |= FiveMinutes

		case "10_minutes":
			parsedResolution |= TenMinutes

		case "hourly":
			parsedResolution |= Hourly

		case "subdaily":
			parsedResolution |= SubDaily

		case "daily":
			parsedResolution |= Daily

		case "monthly":
			parsedResolution |= Monthly

		case "annual":
			parsedResolution |= Annually

		case "multi_annual":
			parsedResolution |= MultiAnnually

		}
	}
	*r = parsedResolution
}

// nolint: goconst
func (r *Resolution) ParseUrlValues(q url.Values, key string) {
	var parsedResolution Resolution
	parts := q[key]
	for _, part := range parts {
		switch part {
		case "1_minute":
			parsedResolution |= OneMinute

		case "5_minutes":
			parsedResolution |= FiveMinutes

		case "10_minutes":
			parsedResolution |= TenMinutes

		case "hourly":
			parsedResolution |= Hourly

		case "subdaily":
			parsedResolution |= SubDaily

		case "daily":
			parsedResolution |= Daily

		case "monthly":
			parsedResolution |= Monthly

		case "annual":
			parsedResolution |= Annually

		case "multi_annual":
			parsedResolution |= MultiAnnually

		}
	}
	*r = parsedResolution
}

func (r Resolution) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Resolution) UnmarshalJSON(src []byte) error {
	var resolution string
	if err := json.Unmarshal(src, &resolution); err != nil {
		return err
	}
	r.ParseString(resolution)
	return nil
}

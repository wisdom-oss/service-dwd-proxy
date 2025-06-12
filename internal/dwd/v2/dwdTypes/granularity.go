package dwdTypes

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

type Granularity uint8

const (
	Granularity_None   Granularity = 0
	Granularity_Annual Granularity = iota
	Granularity_Monthly
	Granularity_Daily
	Granularity_SubDaily
	Granularity_Hourly
	Granularity_Every10Mins
	Granularity_Every5Mins
	Granularity_EveryMinute
)

func (g Granularity) String() string {
	switch g {
	case Granularity_None:
		return ""
	case Granularity_Annual:
		return "annual"
	case Granularity_Monthly:
		return "monthly"
	case Granularity_Daily:
		return "daily"
	case Granularity_SubDaily:
		return "subDaily"
	case Granularity_Hourly:
		return "hourly"
	case Granularity_Every10Mins:
		return "every10Minutes"
	case Granularity_Every5Mins:
		return "every5Minutes"
	case Granularity_EveryMinute:
		return "everyMinute"
	default:
		return ""
	}
}

func (g Granularity) UrlPart() string {
	switch g {
	case Granularity_Annual:
		return g.String()
	case Granularity_Monthly:
		return g.String()
	case Granularity_Daily:
		return g.String()
	case Granularity_SubDaily:
		return "subdaily"
	case Granularity_Hourly:
		return g.String()
	case Granularity_Every10Mins:
		return "10_minutes"
	case Granularity_Every5Mins:
		return "5_minutes"
	case Granularity_EveryMinute:
		return "1_minute"
	default:
		return ""
	}
}

func (g *Granularity) Parse(src any) error {
	if v := reflect.ValueOf(src); !v.IsValid() {
		return errors.New("granularity may not be <nil>")
	}

	switch v := src.(type) {
	case string:
		return g.parseString(v)
	case []byte:
		return g.parseString(string(v))
	default:
		return errors.New("unsupported input type")
	}

}

func (g *Granularity) parseString(s string) error {
	switch strings.TrimSpace(s) {
	case Granularity_Annual.String(), Granularity_Annual.UrlPart():
		*g = Granularity_Annual
	case Granularity_Monthly.String(), Granularity_Monthly.UrlPart():
		*g = Granularity_Monthly
	case Granularity_Daily.String(), Granularity_Daily.UrlPart():
		*g = Granularity_Daily
	case Granularity_SubDaily.String(), Granularity_SubDaily.UrlPart():
		*g = Granularity_SubDaily
	case Granularity_Hourly.String(), Granularity_Hourly.UrlPart():
		*g = Granularity_Hourly
	case Granularity_Every10Mins.String(), Granularity_Every10Mins.UrlPart():
		*g = Granularity_Every10Mins
	case Granularity_Every5Mins.String(), Granularity_Every5Mins.UrlPart():
		*g = Granularity_Every5Mins
	case Granularity_EveryMinute.String(), Granularity_EveryMinute.UrlPart():
		*g = Granularity_EveryMinute
	default:
		*g = Granularity_None
		return errors.New("unsupported granularity")
	}
	return nil
}

func (g Granularity) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

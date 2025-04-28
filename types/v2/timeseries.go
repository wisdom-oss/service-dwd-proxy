package v2

type Timeseries struct {
	Datapoints       []Datapoint     `json:"datapoints"`
	Metadata         []FieldMetadata `json:"metadata"`
	DescriptionFiles []File          `json:"descriptionFiles"`
}

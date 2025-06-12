package v2

type HealthStatus struct {
	Healthy bool   `csv:"healthy"          json:"healthy"`
	Reason  string `csv:"reason,omitempty" json:"reason,omitempty"`
}

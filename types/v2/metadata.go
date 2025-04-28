package v2

import "time"

type FieldMetadata struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Unit        string    `json:"unit"`
	ValidFrom   time.Time `json:"-"`
	ValidUntil  time.Time `json:"-"`
}

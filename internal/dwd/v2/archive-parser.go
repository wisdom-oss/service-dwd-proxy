package v2

import (
	"microservice/internal/dwd/v2/internal/parser"
	v2 "microservice/types/v2"
)

func HandleArchive(filepath string) (datapoints []v2.Datapoint, metadata []v2.FieldMetadata, err error) {
	return parser.ReadAchrive(filepath)
}

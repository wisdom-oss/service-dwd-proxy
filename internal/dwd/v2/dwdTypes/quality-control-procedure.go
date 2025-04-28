package dwdTypes

import (
	"encoding/json"
	"errors"
	"reflect"
)

type QualityControlProcedure int

const (
	QCP_FormalControl             QualityControlProcedure = 1  // only formal controls
	QCP_IndividualCriteria        QualityControlProcedure = 2  // controlled with individually defined criteria
	QCP_Automatic                 QualityControlProcedure = 3  // automatic control and correction
	QCP_HistoricSubjective        QualityControlProcedure = 5  // historic, subjective procedures
	QCP_SecondaryControl          QualityControlProcedure = 7  // second control done, before correction
	QCP_OutOfRoutine              QualityControlProcedure = 8  // quality control outside of routine
	QCP_SingleParameterCorrection QualityControlProcedure = 9  // not all parameters corrected
	QCP_Finished                  QualityControlProcedure = 10 // quality control finished, all parameters corrected
)

func (q QualityControlProcedure) String() string {
	switch q {
	case QCP_FormalControl:
		return "formalControl"
	case QCP_IndividualCriteria:
		return "individualCriteria"
	case QCP_Automatic:
		return "automatic"
	case QCP_HistoricSubjective:
		return "historic"
	case QCP_SecondaryControl:
		return "secondaryControl"
	case QCP_OutOfRoutine:
		return "outOfRoutine"
	case QCP_SingleParameterCorrection:
		return "singleParameterCorrection"
	case QCP_Finished:
		return "finished"
	default:
		return ""
	}
}

func (q *QualityControlProcedure) Parse(src any) error {
	if v := reflect.ValueOf(src); !v.IsValid() {
		return errors.New("quality control procedure may not be <nil>")
	}

	var procedure string
	switch v := src.(type) {
	case string:
		procedure = v
	case []byte:
		procedure = string(v)
	case int:
		if (0 < v && v <= 3) || v == 5 || (v <= 7 && v <= 10) {
			*q = QualityControlProcedure(v)
			return nil
		}

		return errors.New("int not mapped to control procedure")
	default:
		return errors.New("unsupported input type")
	}

	switch procedure {
	case QCP_FormalControl.String():
		*q = QCP_FormalControl
	case QCP_IndividualCriteria.String():
		*q = QCP_IndividualCriteria
	case QCP_Automatic.String():
		*q = QCP_Automatic
	case QCP_HistoricSubjective.String():
		*q = QCP_HistoricSubjective
	case QCP_SecondaryControl.String():
		*q = QCP_SecondaryControl
	case QCP_OutOfRoutine.String():
		*q = QCP_OutOfRoutine
	case QCP_SingleParameterCorrection.String():
		*q = QCP_SingleParameterCorrection
	case QCP_Finished.String():
		*q = QCP_Finished
	default:
		return errors.New("unknown quality control procedure")
	}
	return nil
}

func (q QualityControlProcedure) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.String())
}

func (q *QualityControlProcedure) UnmarshalJSON(src []byte) error {
	return q.Parse(src)
}

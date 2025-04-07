package dwdTypes

import (
	"encoding/json"
	"errors"
	"reflect"
)

type QualityFlag int

const (
	QF_Unflagged                      QualityFlag = 0    // unflagged
	QF_NoObjections                   QualityFlag = 1    // no objections (either checked or unchecked)
	QF_Corrected                      QualityFlag = 2    // corrected
	QF_ConfirmedWithRejectedObjection QualityFlag = 3    // confirmed with rejected objection
	QF_AddedOrCalculated              QualityFlag = 4    // added or calculated
	QF_Objected                       QualityFlag = 5    // objected
	QF_OnlyFormalCheck                QualityFlag = 6    // only formally checked
	QF_FormalObjection                QualityFlag = 7    // formal objection
	QF_FlagNonExisitent               QualityFlag = -999 // quality flag does not exist
)

func (qf QualityFlag) String() string {
	switch qf {
	case QF_Unflagged:
		return "unflagged"
	case QF_NoObjections:
		return "noObjections"
	case QF_Corrected:
		return "corrected"
	case QF_ConfirmedWithRejectedObjection:
		return "confirmedAfterObjection"
	case QF_AddedOrCalculated:
		return "addedOrCalculated"
	case QF_Objected:
		return "objected"
	case QF_OnlyFormalCheck:
		return "onlyFormalCheck"
	case QF_FormalObjection:
		return "formalObjection"
	case QF_FlagNonExisitent:
		return "flagNonExistent"
	default:
		return ""
	}
}

func (qf *QualityFlag) Parse(src any) error {
	if v := reflect.ValueOf(src); !v.IsValid() || v.IsNil() {
		*qf = QF_FlagNonExisitent
		return nil
	}

	var flag string
	switch v := src.(type) {
	case string:
		flag = v
	case []byte:
		flag = string(v)
	case int:
		if 0 <= v && v <= 7 {
			*qf = QualityFlag(v)
			return nil
		}

		*qf = QF_FlagNonExisitent
		return nil

	default:
		return errors.New("unsupported input type")
	}

	switch flag {
	case QF_Unflagged.String():
		*qf = QF_Unflagged
	case QF_NoObjections.String():
		*qf = QF_NoObjections
	case QF_Corrected.String():
		*qf = QF_Corrected
	case QF_ConfirmedWithRejectedObjection.String():
		*qf = QF_ConfirmedWithRejectedObjection
	case QF_AddedOrCalculated.String():
		*qf = QF_AddedOrCalculated
	case QF_Objected.String():
		*qf = QF_Objected
	case QF_OnlyFormalCheck.String():
		*qf = QF_OnlyFormalCheck
	case QF_FormalObjection.String():
		*qf = QF_FormalObjection
	default:
		*qf = QF_FlagNonExisitent
	}
	return nil

}

func (qf QualityFlag) MarshalJSON() ([]byte, error) {
	if qf == QF_FlagNonExisitent {
		return nil, nil
	}
	return json.Marshal(qf.String())
}

func (qf *QualityFlag) UnmarshalJSON(src []byte) error {
	return qf.Parse(src)
}

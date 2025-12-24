package models

import "github.com/civet148/sqlca/v3"

type FrozenState int

const (
	FrozenState_False = 0
	FrozenState_Ture  = 1
)

func (s FrozenState) String() string {
	switch s {
	case FrozenState_Ture:
		return "True"
	case FrozenState_False:
		return "False"
	}
	return "<FrozenState_Unknown>"
}

type ProductExtraData struct {
	AvgPrice   sqlca.Decimal `json:"avg_price"`   //均价
	SpecsValue string        `json:"specs_value"` //规格
}

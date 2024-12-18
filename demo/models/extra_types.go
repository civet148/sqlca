package models

import "github.com/civet148/sqlca/v2"

type ProductExtraData struct {
	AvgPrice   sqlca.Decimal `json:"avg_price"`   //均价
	SpecsValue string        `json:"specs_value"` //规格
}


package sqlca

import "github.com/civet148/decimal"

type Decimal = decimal.Decimal

func NewDecimal(v any, rounds ...int32) Decimal {
	return decimal.NewDecimal(v, rounds...)
}

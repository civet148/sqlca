package sqlca

import (
	"database/sql/driver"
	"fmt"
	"github.com/civet148/log"
	"math"
	"math/big"

	"github.com/shopspring/decimal"
	//"go.mongodb.org/mongo-driver/x/bsonx"
	//"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	fil   = "1000000000000000000"
	ether = "1000000000000000000"
	btc   = "100000000"
)

type Decimal struct {
	dec decimal.Decimal
}

func NewDecimal(v interface{}) (d Decimal) {
	var err error
	switch v.(type) {
	case int8:
		d.dec = decimal.NewFromInt32(int32(v.(int8)))
	case int16:
		d.dec = decimal.NewFromInt32(int32(v.(int16)))
	case int32:
		d.dec = decimal.NewFromInt32(v.(int32))
	case int64:
		d.dec = decimal.NewFromInt(v.(int64))
	case int:
		d.dec = decimal.NewFromInt(int64(v.(int)))
	case uint8:
		d.dec = decimal.NewFromInt32(int32(v.(uint8)))
	case uint16:
		d.dec = decimal.NewFromInt32(int32(v.(uint16)))
	case uint32:
		d.dec = decimal.NewFromInt32(int32(v.(uint32)))
	case uint64:
		d.dec = decimal.NewFromInt(int64(v.(uint64)))
	case uint:
		d.dec = decimal.NewFromInt(int64(v.(uint)))
	case float32:
		d.dec = decimal.NewFromFloat32(v.(float32))
	case float64:
		d.dec = decimal.NewFromFloat(v.(float64))
	case []byte:
		amt := string(v.([]byte))
		if amt == "" {
			amt = "0"
		}
		d.dec, _ = decimal.NewFromString(amt)
	case string:
		amt := v.(string)
		if amt == "" {
			amt = "0"
		}
		d.dec, err = decimal.NewFromString(v.(string))
		if err != nil {
			log.Errorf("value [%v] is not a valid number", v)
		}
	default:
		d.dec, err = decimal.NewFromString(fmt.Sprintf("%v", v))
		if err != nil {
			log.Errorf("value [%v] is not a valid number", v)
		}
	}
	return
}

func (d Decimal) BigInt() (b *big.Int, ok bool) {
	b = new(big.Int)
	return b.SetString(d.String(), 10)
}

func (d Decimal) Amount2FIL() Decimal {
	return d.Mul(NewDecimal(fil))
}

func (d Decimal) FIL2Amount() Decimal {
	return d.Div(NewDecimal(fil))
}

func (d Decimal) Amount2Ether() Decimal {
	return d.Mul(NewDecimal(ether))
}

func (d Decimal) Ether2Amount() Decimal {
	return d.Div(NewDecimal(ether))
}

func (d Decimal) Amount2Btc() Decimal {
	return d.Mul(NewDecimal(btc))
}

func (d Decimal) Btc2Amount() Decimal {
	return d.Div(NewDecimal(btc))
}

func (d Decimal) Amount2Coin(prec int) Decimal {
	if prec < 0 {
		panic("precision cannot be negative")
	}
	return d.Mul(NewDecimal(math.Pow10(prec)))
}

func (d Decimal) Coin2Amount(prec int) Decimal {
	if prec < 0 {
		panic("precision cannot be negative")
	}
	return d.Div(NewDecimal(math.Pow10(prec)))
}

func (d *Decimal) FromString(v string) {
	d.dec, _ = decimal.NewFromString(v)
}

func (d *Decimal) FromFloat(v float64) {
	d.dec = decimal.NewFromFloat(v)
}

func (d *Decimal) FromInt(v int64) {
	d.dec = decimal.NewFromInt(v)
}

func convertDecimal(d interface{}) (dv Decimal) {
	var ok bool
	if dv, ok = d.(Decimal); !ok {
		dv = NewDecimal(d)
	}
	return
}

// Add returns d + d2
func (d Decimal) Add(d2 interface{}) Decimal {
	d3 := convertDecimal(d2)
	return Decimal{
		dec: d.dec.Add(d3.dec),
	}
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
	return Decimal{
		dec: d.dec.Abs(),
	}
}

// Sub returns d - d2.
func (d Decimal) Sub(d2 interface{}) Decimal {
	d3 := convertDecimal(d2)
	return Decimal{
		dec: d.dec.Sub(d3.dec),
	}
}

// Neg returns -d.
func (d Decimal) Neg() Decimal {
	return Decimal{
		dec: d.dec.Neg(),
	}
}

// Mul returns d * d2.
func (d Decimal) Mul(d2 interface{}) Decimal {
	d3 := convertDecimal(d2)
	return Decimal{
		dec: d.dec.Mul(d3.dec),
	}
}

// Div returns d / d2. If it doesn't divide exactly, the result will have
// DivisionPrecision digits after the decimal point.
func (d Decimal) Div(d2 interface{}) Decimal {
	d3 := convertDecimal(d2)
	return Decimal{
		dec: d.dec.Div(d3.dec),
	}
}

// Mod returns d % d2.
func (d Decimal) Mod(d2 interface{}) Decimal {
	d3 := convertDecimal(d2)
	return Decimal{
		dec: d.dec.Mod(d3.dec),
	}
}

// Pow returns d to the power d2
func (d Decimal) Pow(d2 interface{}) Decimal {
	d3 := convertDecimal(d2)
	return Decimal{
		dec: d.dec.Pow(d3.dec),
	}
}

// Cmp compares the numbers represented by d and d2 and returns:
//
//	-1 if d <  d2
//	 0 if d == d2
//	+1 if d >  d2
func (d Decimal) Cmp(d2 interface{}) int {
	d3 := convertDecimal(d2)
	return d.dec.Cmp(d3.dec)
}

// Equal returns whether the numbers represented by d and d2 are equal.
func (d Decimal) Equal(d2 interface{}) bool {
	d3 := convertDecimal(d2)
	return d.dec.Equal(d3.dec)
}

// GreaterThan (GT) returns true when d is greater than d2.
func (d Decimal) GreaterThan(d2 interface{}) bool {
	d3 := convertDecimal(d2)
	return d.dec.GreaterThan(d3.dec)
}

// GreaterThanOrEqual (GTE) returns true when d is greater than or equal to d2.
func (d Decimal) GreaterThanOrEqual(d2 interface{}) bool {
	d3 := convertDecimal(d2)
	return d.dec.GreaterThanOrEqual(d3.dec)
}

// LessThan (LT) returns true when d is less than d2.
func (d Decimal) LessThan(d2 interface{}) bool {
	d3 := convertDecimal(d2)
	return d.dec.LessThan(d3.dec)
}

// LessThanOrEqual (LTE) returns true when d is less than or equal to d2.
func (d Decimal) LessThanOrEqual(d2 interface{}) bool {
	d3 := convertDecimal(d2)
	return d.dec.LessThanOrEqual(d3.dec)
}

// Sign returns:
//
//	-1 if d <  0
//	 0 if d == 0
//	+1 if d >  0
func (d Decimal) Sign() int {
	return d.dec.Sign()
}

// IsPositive return
//
//	true if d > 0
//	false if d == 0
//	false if d < 0
func (d Decimal) IsPositive() bool {
	return d.dec.IsPositive()
}

// IsNegative return
//
//	true if d < 0
//	false if d == 0
//	false if d > 0
func (d Decimal) IsNegative() bool {
	return d.dec.IsNegative()
}

// IsZero return
//
//	true if d == 0
//	false if d > 0
//	false if d < 0
func (d Decimal) IsZero() bool {
	return d.dec.IsZero()
}

// IntPart returns the integer component of the decimal.
func (d Decimal) IntPart() int64 {
	return d.dec.IntPart()
}

// Float64 returns the nearest float64 value for d and a bool indicating
// whether f represents d exactly.
func (d Decimal) Float64() (f float64) {
	f, _ = d.dec.Float64()
	return
}

// String returns the string representation of the decimal
// with the fixed point.
//
// Example:
//
//	d := New(-12345, -3)
//	println(d.String())
//
// Output:
//
//	-12.345
func (d Decimal) String() string {
	return d.dec.String()
}

// StringFixed returns a rounded fixed-point string with places digits after
// the decimal point.
//
// Example:
//
//	NewFromFloat(0).StringFixed(2) // output: "0.00"
//	NewFromFloat(0).StringFixed(0) // output: "0"
//	NewFromFloat(5.45).StringFixed(0) // output: "5"
//	NewFromFloat(5.45).StringFixed(1) // output: "5.5"
//	NewFromFloat(5.45).StringFixed(2) // output: "5.45"
//	NewFromFloat(5.45).StringFixed(3) // output: "5.450"
//	NewFromFloat(545).StringFixed(-1) // output: "550"
func (d Decimal) StringFixed(places int32) string {
	return d.dec.StringFixed(places)
}

// Round rounds the decimal to places decimal places.
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Example:
//
//	NewFromFloat(5.45).Round(1).String() // output: "5.5"
//	NewFromFloat(545).Round(-1).String() // output: "550"
func (d Decimal) Round(places int32) Decimal {

	return Decimal{
		dec: d.dec.Round(places),
	}
}

// Truncate truncates off digits from the number, without rounding.
//
// NOTE: precision is the last digit that will not be truncated (must be >= 0).
//
// Example:
//
//	decimal.NewFromString("123.456").Truncate(2).String() // "123.45"
func (d Decimal) Truncate(precision int32) Decimal {
	return Decimal{
		dec: d.dec.Truncate(precision),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return d.dec.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	return d.dec.UnmarshalJSON(decimalBytes)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Decimal) MarshalBinary() (data []byte, err error) {
	return d.dec.MarshalBinary()
}

// MarshalBSON implements the bson.Marshaler interface.
func (d Decimal) MarshalBSON() ([]byte, error) {
	return d.dec.MarshalJSON()
}

func (d *Decimal) UnmarshalBSON(data []byte) error {
	return d.dec.UnmarshalJSON(data)
}

//
////MarshalBSONValue implements the bson.Marshaler interface.
//func (d Decimal) MarshalBSONValue() (bsontype.Type, []byte, error) {
//	pd, err := primitive.ParseDecimal128(d.dec.String())
//	if err != nil {
//		return bsontype.Decimal128, nil, log.Errorf(err.Error())
//	}
//	return bsonx.Decimal128(pd).MarshalBSONValue()
//}
//
//// UnmarshalBSONValue implements the bson.Unmarshaler interface.
//func (d *Decimal) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
//
//	if bt == bsontype.Decimal128 {
//		bd, _, ok := bsoncore.ReadDecimal128(data)
//		if !ok {
//			return log.Errorf("unmarshal decimal128 error")
//		}
//		d.FromString(bd.String())
//	} else if bt == bsontype.Double {
//		bd, _, ok := bsoncore.ReadDouble(data)
//		if !ok {
//			return log.Errorf("unmarshal Double error")
//		}
//		d.FromFloat(bd)
//	} else if bt == bsontype.String {
//		bd, _, ok := bsoncore.ReadString(data)
//		if !ok {
//			return log.Errorf("unmarshal String error")
//		}
//		d.FromString(bd)
//	} else if bt == bsontype.Int64 {
//		bd, _, ok := bsoncore.ReadInt64(data)
//		if !ok {
//			return log.Errorf("unmarshal Int64 error")
//		}
//		d.FromInt(bd)
//	} else if bt == bsontype.Int32 {
//		bd, _, ok := bsoncore.ReadInt32(data)
//		if !ok {
//			return log.Errorf("unmarshal Int32 error")
//		}
//		d.FromInt(int64(bd))
//	} else {
//		return log.Errorf("unknown bson type [%s] to unmarshal", bt)
//	}
//	return nil
//}
//
//// GetBSON implements the bson.Getter interface (mgo.v2)
//func (d Decimal) GetBSON() (interface{}, error) {
//	return d.dec.String(), nil
//}
//
//// SetBSON implements the bson.Setter interface (mgo.v2)
//func (d *Decimal) SetBSON(raw bson.Raw) error {
//	var err error
//	var strData string
//	if err = raw.Unmarshal(&strData); err != nil {
//		var d128 bson.Decimal128
//		if err = raw.Unmarshal(&d128); err != nil {
//			fmt.Printf("SetBSON call raw.Unmarshal to decimal128 error [%s]\n", err)
//			return err
//		}
//		*d = NewDecimal(d128.String())
//	} else {
//		if err = d.dec.UnmarshalJSON([]byte(strData)); err != nil {
//			fmt.Printf("SetBSON call dec.UnmarshalJSON [%s] error [%s]\n", strData, err)
//			return err
//		}
//	}
//	return nil
//}

func (d Decimal) Marshal() ([]byte, error) {
	return d.dec.MarshalJSON()
}

func (d *Decimal) Unmarshal(data []byte) error {
	return d.dec.UnmarshalJSON(data)
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface. As a string representation
// is already used when encoding to text, this method stores that string as []byte
func (d *Decimal) UnmarshalBinary(data []byte) error {
	return d.dec.UnmarshalBinary(data)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Decimal) Scan(src interface{}) error {
	return d.dec.Scan(src)
}

// Value implements the driver.Valuer interface for database serialization.
func (d Decimal) Value() (driver.Value, error) {
	return d.dec.Value()
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization.
func (d Decimal) MarshalText() (text []byte, err error) {
	return d.dec.MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization.
func (d *Decimal) UnmarshalText(text []byte) error {
	return d.dec.UnmarshalText(text)
}

// StringScaled first scales the decimal then calls .String() on it.
// NOTE: buggy, unintuitive, and DEPRECATED! Use StringFixed instead.
func (d Decimal) StringScaled(exp int32) string {

	return d.dec.StringScaled(exp)
}

// Min returns the smallest Decimal that was passed in the arguments.
// To call this function with an array, you must do:
// This makes it harder to accidentally call Min with 0 arguments.
func (d Decimal) Min(rest ...Decimal) Decimal {

	var r []decimal.Decimal
	for _, v := range rest {
		r = append(r, v.dec)
	}
	return Decimal{dec: decimal.Min(d.dec, r...)}
}

// Max returns the largest Decimal that was passed in the arguments.
// To call this function with an array, you must do:
// This makes it harder to accidentally call Max with 0 arguments.
func (d Decimal) Max(rest ...Decimal) Decimal {
	var r []decimal.Decimal
	for _, v := range rest {
		r = append(r, v.dec)
	}
	return Decimal{dec: decimal.Max(d.dec, r...)}
}

// Sum returns the combined total of the provided first and rest Decimals
func (d Decimal) Sum(rest ...Decimal) Decimal {
	var r []decimal.Decimal
	for _, v := range rest {
		r = append(r, v.dec)
	}
	return Decimal{dec: decimal.Sum(d.dec, r...)}
}

// Sin returns the sine of the radian argument x.
func (d Decimal) Sin() Decimal {
	return Decimal{dec: d.dec.Sin()}
}

// Cos returns the cosine of the radian argument x.
func (d Decimal) Cos() Decimal {
	return Decimal{dec: d.dec.Cos()}
}

// Tan returns the tangent of the radian argument x.
func (d Decimal) Tan() Decimal {
	return Decimal{dec: d.dec.Tan()}
}

// GetDecimal returns the decimal.Decimal type
func (d Decimal) GetDecimal() decimal.Decimal {
	return d.dec
}

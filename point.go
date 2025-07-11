package sqlca

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Point 表示数据库中的point类型
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Value 实现driver.Valuer接口，将Point转换为数据库可存储的格式
func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%f,%f)", p.X, p.Y), nil
}

// Scan 实现sql.Scanner接口，将数据库中的值转换为Point类型
func (p *Point) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("unsupported type for Point.Scan: %T", value)
	}

	// 解析格式如"(123.45,67.89)"的字符串
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid point format: %s", s)
	}

	var err error
	p.X, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return fmt.Errorf("invalid x coordinate: %s, error: %v", parts[0], err)
	}

	p.Y, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return fmt.Errorf("invalid y coordinate: %s, error: %v", parts[1], err)
	}

	return nil
}

// MarshalJSON 实现json.Marshaler接口
func (p Point) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]float64{
		"x": p.X,
		"y": p.Y,
	})
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (p *Point) UnmarshalJSON(data []byte) error {
	var m map[string]float64
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if x, ok := m["x"]; ok {
		p.X = x
	}

	if y, ok := m["y"]; ok {
		p.Y = y
	}

	return nil
}

// String 提供Point的字符串表示
func (p Point) String() string {
	return fmt.Sprintf("POINT(%f,%f)", p.X, p.Y)
}

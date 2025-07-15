package sqlca

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/civet148/sqlca/v3/types"
	"math"
)

// Point 表示数据库中的point类型
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Value 实现driver.Valuer接口，将Point转换为数据库可存储的格式
func (p Point) Value() (driver.Value, error) {
	return types.NewSqlClauseValue("POINT(%f,%f)", p.X, p.Y), nil
}

// Scan 实现sql.Scanner接口，用于解析MySQL的POINT类型数据(0x000000000101000000E5D022DBF98E5B40F6285C8FC27524C0 => POINT(110.234 -10.23))
func (p *Point) Scan(src any) error {
	// 检查输入类型
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("invalid type for Point: %T", src)
	}

	hexStr := fmt.Sprintf("0x%X", data)
	data = []byte(hexStr)
	// 检查数据是否为空
	if len(data) == 0 {
		return nil
	}

	// 处理十六进制格式的数据
	// 移除可能的前缀"0x"
	data = bytes.TrimPrefix(data, []byte("0x"))

	// 检查数据长度是否足够
	if len(data) < 42 { // 至少需要42个十六进制字符 (21个字节)
		return fmt.Errorf("invalid Point data: too short, length: %d", len(data))
	}

	// 解析WKB格式的POINT数据
	// 前9个字节是头部信息，后面16个字节是坐标数据
	// 这里我们假设数据是小端字节序

	// 解析X坐标 (第10-17字节)
	xBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		_, err := fmt.Sscanf(string(data[18+2*i:20+2*i]), "%02x", &xBytes[i])
		if err != nil {
			return fmt.Errorf("parse X coordinate failed: %w", err)
		}
	}
	p.X = math.Float64frombits(binary.LittleEndian.Uint64(xBytes))

	// 解析Y坐标 (第18-25字节)
	yBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		_, err := fmt.Sscanf(string(data[34+2*i:36+2*i]), "%02x", &yBytes[i])
		if err != nil {
			return fmt.Errorf("parse Y coordinate failed: %w", err)
		}
	}
	p.Y = math.Float64frombits(binary.LittleEndian.Uint64(yBytes))

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

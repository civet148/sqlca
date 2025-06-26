package sqlca

import (
	"database/sql/driver"
	"fmt"
	"testing"
)

type object struct {
	Name   string `json:"name"`
	Gender int32  `json:"gender"`
}

type objectValuer struct {
	Name   string `json:"name"`
	Gender int32  `json:"gender"`
}

func (o objectValuer) Value() (driver.Value, error) {
	return string("hello"), nil
}

func TestIndirectValue(t *testing.T) {
	var obj *object
	fmt.Printf("object nil => [%v]\n", indirectValue(obj))
	fmt.Printf("object => [%v]\n", indirectValue(object{
		Name:   "lory1",
		Gender: 1,
	}))
	fmt.Printf("object pointer => [%v]\n", indirectValue(&object{
		Name:   "lory2",
		Gender: 1,
	}))
	fmt.Printf("object valuer => [%v]\n", indirectValue(&objectValuer{
		Name:   "lory3",
		Gender: 1,
	}))
	fmt.Printf("map => [%v]\n", indirectValue(map[string]any{
		"name": "goo",
		"age":  25,
	}))
	fmt.Printf("slice => [%v]\n", indirectValue([]any{
		"goo",
		25,
	}))
}

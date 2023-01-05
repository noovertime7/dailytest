package utils

import (
	"encoding/json"
	"fmt"
)

func ParseToJsonByte(obj string) ([]byte, error) {
	fmt.Println(obj)
	return json.Marshal(obj)
}

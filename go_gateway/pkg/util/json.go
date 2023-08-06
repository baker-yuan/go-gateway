package util

import "encoding/json"

// Unmarshal 反序列化
func Unmarshal[T any](data []byte) (T, error) {
	var t T
	return t, json.Unmarshal(data, &t)
}

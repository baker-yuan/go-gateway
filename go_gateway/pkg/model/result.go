package model

import "encoding/json"

type Result struct {
	Code    uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *Result) ToString() string {
	marshal, _ := json.Marshal(r)
	return string(marshal)
}

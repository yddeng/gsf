package json

import (
	"encoding/json"
)

/*
 json 协议解析
*/

type Json struct{}

func (_ Json) Unmarshal(data []byte, o interface{}) (err error) {
	return json.Unmarshal(data, o)
}

func (_ Json) Marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

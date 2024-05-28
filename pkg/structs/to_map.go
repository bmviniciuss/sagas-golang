package structs

import (
	"encoding/json"

	"github.com/bytedance/sonic"
)

func ToMap(data interface{}) (map[string]interface{}, error) {
	b, err := sonic.Marshal(data)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func FromBytes(data []byte, dest interface{}) error {
	err := sonic.Unmarshal(data, dest)
	if err != nil {
		return err
	}
	return nil
}

func ToBytes(data interface{}) ([]byte, error) {
	return sonic.Marshal(data)
}

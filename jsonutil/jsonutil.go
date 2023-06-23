package jsonutil

import "encoding/json"

func Adapt(ori, target interface{}) error {
	bytes, err := json.Marshal(ori)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func MarshalToString(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

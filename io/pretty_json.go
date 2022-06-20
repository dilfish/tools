package io

import (
	"encoding/json"

	"github.com/TylerBrock/colorjson"
)

// PrettyJson is a stupid function
// you should first self-define a formatter
// it returned "" when error but no error
func PrettyJson(obj interface{}) (string, error) {
	origJs, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(origJs, &resp)
	if err != nil {
		return "", err
	}
	str, err := colorjson.Marshal(obj)
	return string(str), err
}

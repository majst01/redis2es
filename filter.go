package main

import (
	"fmt"
	"strings"

	"github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func filter(input string) (string, string, error) {
	contractName := "catchall"
	data := make(map[string]interface{})
	err := json.UnmarshalFromString(input, &data)
	if err != nil {
		return "", contractName, fmt.Errorf("cannot decode data:%v", err)
	}

	// check if contract in any case is present, lowercase then
	for k, v := range data {
		if strings.ToLower(k) == "contract" {
			oldValue := strings.ToLower(v.(string))
			delete(data, k)
			data["contract"] = oldValue
			contractName = oldValue
		}
	}

	result, err := json.MarshalToString(&data)
	if err != nil {
		return "", contractName, fmt.Errorf("cannot encode data:%v", err)
	}
	return result, contractName, nil

}

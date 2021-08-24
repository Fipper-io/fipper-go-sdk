package fipper_go_sdk

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
)

const (
	Rarely     = 15
	Normal     = 7
	Frequently = 3
)

type Flag struct {
	State bool
	Type  int
	Value interface{}
}

type ConfigManager struct {
	Flags map[string]Flag
}

func CreateConfigManagerFromRawData(rawData string) *ConfigManager {
	decodedBase64, _ := base64.StdEncoding.DecodeString(rawData)
	rawJsonData, _ := gzip.NewReader(bytes.NewReader(decodedBase64))
	result, _ := ioutil.ReadAll(rawJsonData)
	var jsonData map[string]interface{}

	if err := json.Unmarshal(result, &jsonData); err != nil {
		panic(err)
	}

	newConfigManager := ConfigManager{Flags: make(map[string]Flag)}

	for slug, element := range jsonData {
		item := element.(map[string]interface{})
		newConfigManager.Flags[slug] = Flag{State: item["state"].(bool), Type: int(item["type"].(float64)), Value: item["value"]}
	}

	return &newConfigManager
}

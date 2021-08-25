package fipper_go_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var DOMAIN = "https://sync2.fipper.io"

type FipperClient struct {
	Rate             int8
	previousSyncDate *time.Time
	config           *ConfigManager
	eTag             *string
}

func (fc *FipperClient) getActualConfig() *ConfigManager {
	var now = time.Now()

	if fc.previousSyncDate != nil && fc.config != nil {
		if now.Sub(*fc.previousSyncDate) < time.Duration(fc.Rate)*time.Second {
			return fc.config
		}
	}

	return nil
}

func getHashUrl(apiToken string, projectId int, eTag string) string {
	return fmt.Sprintf("%v/hash?apiToken=%v&item=%v&eTag=%v", DOMAIN, apiToken, projectId, eTag)
}

func getConfigUrl(apiToken string, projectId int) string {
	return fmt.Sprintf("%v/config?apiToken=%v&item=%v", DOMAIN, apiToken, projectId)
}

func (fc *FipperClient) GetConfig(environment string, apiToken string, projectId int) (*ConfigManager, error) {
	if prevConfig := fc.getActualConfig(); prevConfig != nil {
		return prevConfig, nil
	}

	if fc.previousSyncDate != nil && fc.config != nil && fc.eTag != nil {
		resp, err := http.Head(getHashUrl(apiToken, projectId, *fc.eTag))

		if err != nil {
			if fc.config != nil {
				return fc.config, nil
			}
			return nil, errors.New("can't fetch config hash")
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotModified {
			return fc.config, nil
		}
	}

	resp, err := http.Get(getConfigUrl(apiToken, projectId))

	if err != nil {
		return nil, errors.New("can't fetch config data")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var rawData map[string]interface{}
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			return nil, errors.New("can't fetch config data body")
		}

		if err := json.Unmarshal(body, &rawData); err != nil {
			return nil, errors.New("can't parse config body")
		}

		newEtag := rawData["eTag"].(string)
		newRawConfig := rawData["config"].(map[string]interface{})
		syncDate := time.Now()

		fc.config = CreateConfigManagerFromRawData(newRawConfig[environment].(string))
		fc.eTag = &newEtag
		fc.previousSyncDate = &syncDate
	} else {
		if fc.config == nil {
			return nil, errors.New(fmt.Sprintf("wrong fetch status: %v", resp.StatusCode))
		}
	}

	return fc.config, nil
}

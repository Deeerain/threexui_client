package threexuiclient

import (
	"encoding/json"
	"fmt"
)

type XUIInbound struct {
	ID             int    `json:"id"`
	Remark         string `json:"remark"`
	Enable         bool   `json:"enable"`
	Listen         string `json:"listen"`
	Port           int    `json:"port"`
	Settings       string `json:"settings"`
	StreamSettings string `json:"streamSettings"`
}

func (i *XUIInbound) GetSettings() (*XUIInboundSettings, error) {
	var settings XUIInboundSettings

	if err := json.Unmarshal([]byte(i.Settings), &settings); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &settings, nil

}

func (i *XUIInbound) GetClients() ([]XUIInboundClient, error) {
	settings, err := i.GetSettings()
	if err != nil {
		return nil, fmt.Errorf("get clinets error: %w", err)
	}

	return settings.Cleints, err
}

type XUIInboundSettings struct {
	Cleints []XUIInboundClient `json:"clients"`
}

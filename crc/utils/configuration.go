package utils

import (
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type crcConfig struct {
	BaseUrl      *string `hcl:"base_url"`
	TokenURL     *string `hcl:"token_url"`
	ClientID     *string `hcl:"client_id"`
	ClientSecret *string `hcl:"client_secret"`
}

func ConfigInstance() interface{} {
	return &crcConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) crcConfig {
	if connection == nil || connection.Config == nil {
		return crcConfig{}
	}
	config, ok := connection.Config.(crcConfig)
	if !ok {
		panic(fmt.Errorf("configuration could not be loaded: %v", connection.Config))
	}
	return config
}

package config

import (
	"encoding/json"
	"os"
)

const (
	configFileName = "config.json"
)

// Config Storing all the configurations inside this variable
type Config struct {
	Token   string `json:"token"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	SSLMode string `json:"ssl_mode"`
}

// LoadConfiguration method opens the 'config.json' file
// and imports the server configs into the config variable.
func LoadConfiguration() (*Config, error) {
	var config Config

	configFile, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}

	defer configFile.Close()

	if er := json.NewDecoder(configFile).Decode(&config); er != nil {
		return nil, er
	}

	return &config, nil
}

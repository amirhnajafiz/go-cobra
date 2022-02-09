package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var file = "config.json"

// LoadConfiguration method opens the 'config.json' file
// and imports the server configs into the config variable.
func LoadConfiguration() Config {
	var config Config

	configFile, err := os.Open(file)
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			panic(err.Error())
		}
	}(configFile)

	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	_ = jsonParser.Decode(&config)

	return config
}

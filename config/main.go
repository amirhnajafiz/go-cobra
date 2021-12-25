package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var config = LoadConfiguration("config.json")

type Config struct {
	Token   string `json:"token"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	SSLMode string `json:"ssl_mode"`
}

func LoadConfiguration(file string) Config {
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

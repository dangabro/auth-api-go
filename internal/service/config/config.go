package config

import (
	"auth-api-go/internal/data"
	"encoding/json"
	"fmt"
	"os"
)

func LoadConfig() (error, data.Config) {
	// get environment variable contents
	fileName := os.Getenv(data.ConfigFile)
	if fileName == "" {
		return fmt.Errorf("environment variable %s not set", data.ConfigFile), data.Config{}
	}

	_, err := os.Stat(fileName)
	if err != nil {
		return err, data.Config{}
	}

	// read the file contents and marshal to json
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return err, data.Config{}
	}

	// check return values
	config := data.Config{}
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return err, data.Config{}
	}

	return nil, config
}

package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	var configPath string
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	} else {
		configPath = home + "/" + configFileName
	}
	return configPath, nil
}

func Read() Config {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}
	}
	rawData, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}
	}
	jsonData := Config{}
	if err := json.Unmarshal(rawData, &jsonData); err != nil {
		return Config{}
	}
	return jsonData
}

func write(conf Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	ok := os.WriteFile(configPath, data, 0666)
	if ok != nil {
		return ok
	}
	return nil
}

func (conf *Config) SetUser(user string) error {
	conf.CurrentUserName = user
	if err := write(*conf); err != nil {
		return err
	}
	return nil
}

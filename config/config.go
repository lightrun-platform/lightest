package config

import (
	"encoding/json"
	"os"
)

type AgentConfig struct {
	ApiKey   string `json:"apiKey"`
	Version  string `json:"version"`
	HostName string `json:"hostName"`
}

type Config struct {
	AgentConfig  `json:"agent"`
	ServerUrl    string `json:"serverUrl"`
	CompanyId    string `json:"companyId"`
	UserEmail    string `json:"userEmail"`
	UserPassword string `json:"userPassword"`
	Certificate  string `json:"pinnedCerts"`
}

var configuration Config

func LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&configuration); err != nil {
		return err
	}

	return nil
}

func GetConfig() *Config {
	return &configuration
}

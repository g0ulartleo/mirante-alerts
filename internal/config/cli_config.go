package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type CLIConfig struct {
	APIHost string `json:"api_host"`
	APIKey  string `json:"api_key"`
}

func GetCLIConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".mirante-alerts")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "cli_config.json"), nil
}

func LoadCLIConfig() (*CLIConfig, error) {
	configPath, err := GetCLIConfigPath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &CLIConfig{
			APIHost: "http://127.0.0.1:40169",
			APIKey:  "",
		}, nil
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config CLIConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func SaveCLIConfig(config *CLIConfig) error {
	configPath, err := GetCLIConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

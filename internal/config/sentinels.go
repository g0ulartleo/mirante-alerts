package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/builtin"
	"gopkg.in/yaml.v3"
)

var (
	SentinelConfigs map[string]*sentinel.SentinelConfig
)

func loadConfig(path string) (*sentinel.SentinelConfig, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read yml file: %v", err)
	}
	var sentinelConfig *sentinel.SentinelConfig
	err = yaml.Unmarshal(yamlFile, &sentinelConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yml file: %v", err)
	}
	if sentinelConfig.Interval == "" && sentinelConfig.Cron == "" {
		return nil, fmt.Errorf("misconfiguration for sentinel %s: interval or cron is required", sentinelConfig.ID)
	}
	if sentinelConfig.Interval != "" && sentinelConfig.Cron != "" {
		return nil, fmt.Errorf("misconfiguration for sentinel %s: interval and cron cannot both be set", sentinelConfig.ID)
	}
	if sentinelConfig.Interval != "" {
		interval, err := time.ParseDuration(sentinelConfig.Interval)
		if err != nil {
			return nil, fmt.Errorf("misconfiguration for sentinel %s: failed to parse interval: %v", sentinelConfig.ID, err)
		}
		sentinelConfig.Cron = fmt.Sprintf("@every %s", interval)
	}
	sentinelConfig.Path = strings.Split(path, "/")[1 : len(strings.Split(path, "/"))-1]
	return sentinelConfig, nil
}

func GetSentinelConfig(id string) (*sentinel.SentinelConfig, error) {
	config, ok := SentinelConfigs[id]
	if !ok {
		return nil, fmt.Errorf("sentinel config not found for id: %s", id)
	}
	return config, nil
}

func InitSentinelConfigs() error {
	SentinelConfigs = make(map[string]*sentinel.SentinelConfig)

	if _, err := os.Stat("sentinels"); os.IsNotExist(err) {
		return fmt.Errorf("sentinels directory does not exist")
	}

	err := filepath.Walk("sentinels", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk sentinels: %v", err)
		}
		if !info.IsDir() && strings.HasSuffix(path, ".yml") {
			config, err := loadConfig(path)
			if err != nil {
				return fmt.Errorf("failed to load config from %s: %v", path, err)
			}
			log.Printf("loaded config id %s from path %s", config.ID, path)
			SentinelConfigs[config.ID] = config
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to load sentinel configs: %v", err)
	}
	return nil
}

func InitSentinelFactory() {
	builtin.Register(sentinel.Factory)
}

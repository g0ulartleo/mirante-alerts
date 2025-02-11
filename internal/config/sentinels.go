package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/g0ulartleo/mirante-alerts/internal/sentinel"
	"github.com/g0ulartleo/mirante-alerts/internal/sentinel/builtin"
	"gopkg.in/yaml.v3"
)

var (
	SentinelConfigs map[int]*sentinel.SentinelConfig
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
	return sentinelConfig, nil
}

func GetSentinelConfig(id int) (*sentinel.SentinelConfig, error) {
	config, ok := SentinelConfigs[id]
	if !ok {
		return nil, fmt.Errorf("sentinel config not found for id: %d", id)
	}
	return config, nil
}

func InitSentinelConfigs() {
	SentinelConfigs = make(map[int]*sentinel.SentinelConfig)

	if _, err := os.Stat("sentinels"); os.IsNotExist(err) {
		log.Fatalf("sentinels directory does not exist")
	}

	err := filepath.Walk("sentinels", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("failed to walk sentinels: %v", err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".yml") {
			config, err := loadConfig(path)
			if err != nil {
				log.Printf("failed to load config from %s: %v", path, err)
				return fmt.Errorf("failed to load config from %s: %v", path, err)
			}
			log.Printf("loaded config id %d from path %s", config.ID, path)
			SentinelConfigs[config.ID] = config
		}
		return nil
	})

	if err != nil {
		log.Fatalf("failed to load sentinel configs: %v", err)
	}
}

func InitSentinelFactory() {
	builtin.Register(sentinel.Factory)
}

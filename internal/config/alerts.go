package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/g0ulartleo/mirante-alerts/internal/alert"
	"gopkg.in/yaml.v3"
)

var (
	Alerts map[string]*alert.Alert
)

func loadAlertConfig(path string) (*alert.Alert, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read yml file: %v", err)
	}
	var alert *alert.Alert
	err = yaml.Unmarshal(yamlFile, &alert)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yml file: %v", err)
	}
	if alert.Interval == "" && alert.Cron == "" {
		return nil, fmt.Errorf("misconfiguration for alert %s: interval or cron is required", alert.ID)
	}
	if alert.Interval != "" && alert.Cron != "" {
		return nil, fmt.Errorf("misconfiguration for alert %s: interval and cron cannot both be set", alert.ID)
	}
	if alert.Interval != "" {
		interval, err := time.ParseDuration(alert.Interval)
		if err != nil {
			return nil, fmt.Errorf("misconfiguration for alert %s: failed to parse interval: %v", alert.ID, err)
		}
		alert.Cron = fmt.Sprintf("@every %s", interval)
	}
	alert.Path = strings.Split(path, "/")[1 : len(strings.Split(path, "/"))-1]
	return alert, nil
}

func GetAlertConfig(id string) (*alert.Alert, error) {
	config, ok := Alerts[id]
	if !ok {
		return nil, fmt.Errorf("alert config not found for id: %s", id)
	}
	return config, nil
}

func InitAlerts() error {
	Alerts = make(map[string]*alert.Alert)

	if _, err := os.Stat("alerts"); os.IsNotExist(err) {
		return fmt.Errorf("alerts directory does not exist")
	}

	err := filepath.Walk("alerts", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk alerts: %v", err)
		}
		if !info.IsDir() && strings.HasSuffix(path, ".yml") {
			config, err := loadAlertConfig(path)
			if err != nil {
				return fmt.Errorf("failed to load config from %s: %v", path, err)
			}
			log.Printf("loaded alert id %s from path %s", config.ID, path)
			Alerts[config.ID] = config
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to load alerts: %v", err)
	}
	return nil
}

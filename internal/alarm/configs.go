package alarm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	Alarms map[string]*Alarm
)

func GetAlarmConfig(id string) (*Alarm, error) {
	config, ok := Alarms[id]
	if !ok {
		return nil, fmt.Errorf("alarm config not found for id: %s", id)
	}
	return config, nil
}

func InitAlarms() error {
	Alarms = make(map[string]*Alarm)

	if _, err := os.Stat("alarms"); os.IsNotExist(err) {
		return fmt.Errorf("alarms directory does not exist")
	}

	err := filepath.Walk("alarms", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk alarms: %v", err)
		}
		if !info.IsDir() && strings.HasSuffix(path, ".yml") {
			config, err := loadAlarmConfig(path)
			if err != nil {
				return fmt.Errorf("failed to load config from %s: %v", path, err)
			}
			log.Printf("loaded alarm id %s from path %s", config.ID, path)
			Alarms[config.ID] = config
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to load alarms: %v", err)
	}
	return nil
}

func loadAlarmConfig(path string) (*Alarm, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		// lookup difference between %w and %v
		return nil, fmt.Errorf("failed to read yml file: %w", err)
	}
	var alarm *Alarm
	err = yaml.Unmarshal(yamlFile, &alarm)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yml file: %v", err)
	}
	if alarm.Interval == "" && alarm.Cron == "" {
		return nil, fmt.Errorf("misconfiguration for alarm %s: interval or cron is required", alarm.ID)
	}
	if alarm.Interval != "" && alarm.Cron != "" {
		return nil, fmt.Errorf("misconfiguration for alarm %s: interval and cron cannot both be set", alarm.ID)
	}
	if alarm.Interval != "" {
		interval, err := time.ParseDuration(alarm.Interval)
		if err != nil {
			return nil, fmt.Errorf("misconfiguration for alarm %s: failed to parse interval: %v", alarm.ID, err)
		}
		alarm.Cron = fmt.Sprintf("@every %s", interval)
	}
	alarm.Path = strings.Split(path, "/")[1 : len(strings.Split(path, "/"))-1]
	return alarm, nil
}

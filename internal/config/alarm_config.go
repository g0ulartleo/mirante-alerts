package config

import (
    "fmt"
    "os"
    "strings"
    "time"

    "gopkg.in/yaml.v3"
    "github.com/g0ulartleo/mirante-alerts/internal/alarm"
)

func LoadAlarmConfig(path string) (*alarm.Alarm, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read yml file: %w", err)
	}
	var al *alarm.Alarm
	err = yaml.Unmarshal(yamlFile, &al)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yml file: %w", err)
	}
	if al.Interval == "" && al.Cron == "" {
		return nil, fmt.Errorf("misconfiguration for alarm %s: interval or cron is required", al.ID)
	}
	if al.Interval != "" && al.Cron != "" {
		return nil, fmt.Errorf("misconfiguration for alarm %s: interval and cron cannot both be set", al.ID)
	}
	if al.Interval != "" {
		interval, err := time.ParseDuration(al.Interval)
		if err != nil {
			return nil, fmt.Errorf("misconfiguration for alarm %s: failed to parse interval: %w", al.ID, err)
		}
		al.Cron = fmt.Sprintf("@every %s", interval)
	}
	al.Path = strings.Split(path, "/")[2 : len(strings.Split(path, "/"))-1]
	return al, nil
}

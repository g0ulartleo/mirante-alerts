package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type SignalsDatabaseConfig struct {
	Driver string      `yaml:"driver"`
	MySQL  MySQLConfig `yaml:"mysql,omitempty"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func LoadSignalsDatabaseConfigFromEnv() *SignalsDatabaseConfig {
	InitEnvLoader()
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		log.Fatalf("DB_DRIVER environment variable is required")
	}

	config := &SignalsDatabaseConfig{
		Driver: driver,
	}

	switch driver {
	case "mysql":
		var port int
		if portStr := os.Getenv("DB_PORT"); portStr != "" {
			p, err := strconv.Atoi(portStr)
			if err != nil {
				log.Fatalf("invalid DB_PORT: %v", err)
			}
			port = p
		} else {
			port = 3306
		}

		config.MySQL = MySQLConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		}
	default:
		log.Fatalf("unsupported driver: %s", driver)
	}

	if err := validateConfig(config); err != nil {
		log.Fatalf("invalid signals database config: %v", err)
	}
	return config
}

func validateConfig(config *SignalsDatabaseConfig) error {
	switch config.Driver {
	case "mysql":
		if config.MySQL.Host == "" {
			return fmt.Errorf("mysql host is required")
		}
		if config.MySQL.Port == 0 {
			config.MySQL.Port = 3306
		}
		if config.MySQL.User == "" {
			return fmt.Errorf("mysql user is required")
		}
		if config.MySQL.Password == "" {
			return fmt.Errorf("mysql password is required")
		}
	default:
		return fmt.Errorf("unsupported driver: %s", config.Driver)
	}
	return nil
}

package config

import (
	"os"
	"sync"
)

type Environment struct {
	DBDriver   string
	DBPort     string
	DBHost     string
	DBUser     string
	DBPassword string
	RedisAddr  string
	HTTPPort   string
	HTTPAddr   string
}

var (
	env  *Environment
	once sync.Once
)

func init() {
	InitEnvLoader()
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Env() *Environment {
	once.Do(func() {
		env = &Environment{
			DBDriver:   os.Getenv("DB_DRIVER"),
			DBPort:     os.Getenv("DB_PORT"),
			DBHost:     os.Getenv("DB_HOST"),
			DBUser:     os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			RedisAddr:  getEnvOrDefault("REDIS_ADDR", "127.0.0.1:6379"),
			HTTPPort:   getEnvOrDefault("HTTP_PORT", "40169"),
			HTTPAddr:   getEnvOrDefault("HTTP_ADDR", "127.0.0.1"),
		}
	})

	return env
}

package config

import (
	"os"
	"sync"
)

type Environment struct {
	RedisAddr string
}

var (
	env  *Environment
	once sync.Once
)

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
			RedisAddr: getEnvOrDefault("REDIS_ADDR", "127.0.0.1:6379"),
		}
	})

	return env
}

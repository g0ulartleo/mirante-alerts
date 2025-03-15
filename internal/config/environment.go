package config

import (
	"os"
	"sync"
)

type Environment struct {
	DBDriver     string
	MySQLDBPort       string
	MySQLDBHost       string
	MySQLDBUser       string
	MySQLDBPassword   string
	RedisAddr    string
	HTTPPort     string
	HTTPAddr     string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
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
			DBDriver:     os.Getenv("DB_DRIVER"),
			MySQLDBHost:       os.Getenv("MYSQL_DB_HOST"),
			MySQLDBPort:       os.Getenv("MYSQL_DB_PORT"),
			MySQLDBUser:       os.Getenv("MYSQL_DB_USER"),
			MySQLDBPassword:   os.Getenv("MYSQL_DB_PASSWORD"),
			RedisAddr:    getEnvOrDefault("REDIS_ADDR", "127.0.0.1:6379"),
			HTTPPort:     getEnvOrDefault("HTTP_PORT", "40169"),
			HTTPAddr:     getEnvOrDefault("HTTP_ADDR", "127.0.0.1"),
			SMTPHost:     os.Getenv("SMTP_HOST"),
			SMTPPort:     os.Getenv("SMTP_PORT"),
			SMTPUser:     os.Getenv("SMTP_USER"),
			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		}
	})

	return env
}

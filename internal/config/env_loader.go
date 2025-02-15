//go:build dev

package config

import (
	"github.com/joho/godotenv"
)

func InitEnvLoader() {
	_ = godotenv.Load()
}

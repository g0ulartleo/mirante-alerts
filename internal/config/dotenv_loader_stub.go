//go:build !dev

package config

func InitEnvLoader() {
	// No .env file is loaded in non-development builds.
}

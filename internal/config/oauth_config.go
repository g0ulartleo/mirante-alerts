package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type OAuthConfig struct {
	Enabled        bool     `yaml:"enabled"`
	Provider       string   `yaml:"provider"`
	RedirectURL    string   `yaml:"redirect_url"`
	AllowedDomains []string `yaml:"allowed_domains"`
	AllowedEmails  []string `yaml:"allowed_emails"`
	SessionTimeout string   `yaml:"session_timeout"`
}

type AuthConfig struct {
	OAuth  OAuthConfig `yaml:"oauth"`
	APIKey string      `yaml:"api_key,omitempty"`
}

func GetAuthConfigPath() string {
	return "config/auth.yaml"
}

func LoadAuthConfig() (*AuthConfig, error) {
	configPath := GetAuthConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &AuthConfig{
			OAuth: OAuthConfig{
				Enabled:        false,
				Provider:       "google",
				SessionTimeout: "24h",
			},
			APIKey: Env().APIKey,
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth config file: %w", err)
	}

	var config AuthConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse auth config file: %w", err)
	}
    config.APIKey = Env().APIKey

	if err := validateAuthConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid auth config: %w", err)
	}

	return &config, nil
}

func SaveAuthConfig(config *AuthConfig) error {
	if err := os.MkdirAll("config", 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := "config/auth.yaml"
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal auth config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write auth config file: %w", err)
	}

	return nil
}

func validateAuthConfig(config *AuthConfig) error {
	if !config.OAuth.Enabled {
		return nil
	}

	if os.Getenv("OAUTH_CLIENT_ID") == "" {
		return fmt.Errorf("OAUTH_CLIENT_ID environment variable is required when OAuth is enabled")
	}

	if os.Getenv("OAUTH_CLIENT_SECRET") == "" {
		return fmt.Errorf("OAUTH_CLIENT_SECRET environment variable is required when OAuth is enabled")
	}

	if os.Getenv("OAUTH_JWT_SECRET") == "" {
		return fmt.Errorf("OAUTH_JWT_SECRET environment variable is required when OAuth is enabled")
	}

	if config.OAuth.RedirectURL == "" {
		return fmt.Errorf("oauth redirect_url is required when OAuth is enabled")
	}

	if config.OAuth.Provider != "google" && config.OAuth.Provider != "github" {
		return fmt.Errorf("unsupported oauth provider: %s (supported: google, github)", config.OAuth.Provider)
	}

	if len(config.OAuth.AllowedDomains) == 0 && len(config.OAuth.AllowedEmails) == 0 {
		return fmt.Errorf("either allowed_domains or allowed_emails must be specified")
	}

	for _, domain := range config.OAuth.AllowedDomains {
		if !strings.HasPrefix(domain, "@") {
			return fmt.Errorf("domain '%s' must start with '@'", domain)
		}
	}

	return nil
}

func (c *OAuthConfig) IsEmailAllowed(email string) bool {
	for _, allowedEmail := range c.AllowedEmails {
		if strings.EqualFold(email, allowedEmail) {
			return true
		}
	}

	for _, domain := range c.AllowedDomains {
		if strings.HasSuffix(strings.ToLower(email), strings.ToLower(domain)) {
			return true
		}
	}

	return false
}

func CreateSampleAuthConfig() error {
	sampleConfig := &AuthConfig{
		OAuth: OAuthConfig{
			Enabled:        true,
			Provider:       "google",
			RedirectURL:    "http://localhost:40169/auth/callback",
			AllowedDomains: []string{"@yourcompany.com"},
			AllowedEmails:  []string{"admin@yourcompany.com", "developer@yourcompany.com"},
			SessionTimeout: "24h",
		},
	}

	return SaveAuthConfig(sampleConfig)
}

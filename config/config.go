// config/config.go
package config

import (
	"fmt"
	"os"
)

type Config struct {
	SmtpHost               string
	SmtpPort               string
	SmtpUser               string
	SmtpPass               string
	FromEmail              string
	FromName               string
	ToEmail                string
	ServerHostname         string
	SmtpInsecureSkipVerify bool
}

func LoadConfig() (*Config, error) {
	skipVerify := false
	if os.Getenv("SMTP_INSECURE_SKIP_VERIFY") == "true" {
		skipVerify = true
	}
	cfg := &Config{
		SmtpHost:               os.Getenv("SMTP_HOST"),
		SmtpPort:               os.Getenv("SMTP_PORT"),
		SmtpUser:               os.Getenv("SMTP_USER"),
		SmtpPass:               os.Getenv("SMTP_PASS"),
		FromEmail:              os.Getenv("FROM_EMAIL"),
		FromName:               os.Getenv("FROM_NAME"),
		ToEmail:                os.Getenv("TO_EMAIL"),
		ServerHostname:         os.Getenv("SERVER_HOSTNAME"),
		SmtpInsecureSkipVerify: skipVerify,
	}

	// Validate required configuration
	if cfg.SmtpHost == "" {
		return nil, fmt.Errorf("SMTP_HOST environment variable not set")
	}
	if cfg.SmtpPort == "" {
		return nil, fmt.Errorf("SMTP_PORT environment variable not set")
	}
	if cfg.SmtpUser == "" {
		return nil, fmt.Errorf("SMTP_USER environment variable not set")
	}
	if cfg.SmtpPass == "" {
		return nil, fmt.Errorf("SMTP_PASS environment variable not set")
	}
	if cfg.FromEmail == "" {
		return nil, fmt.Errorf("FROM_EMAIL environment variable not set")
	}
	if cfg.ToEmail == "" {
		return nil, fmt.Errorf("TO_EMAIL environment variable not set")
	}

	return cfg, nil
}

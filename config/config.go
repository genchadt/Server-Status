// config/config.go
package config

import (
	"fmt"
	"os"
	"time"
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
	DailyReportTime        time.Duration
}

func LoadConfig() (*Config, error) {
	dailyReportTimeStr := os.Getenv("DAILY_REPORT_TIME")
	if dailyReportTimeStr == "" {
		dailyReportTimeStr = "08:00" // Default to 8 AM
	}

	dailyReportTime, err := parseTime(dailyReportTimeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DAILY_REPORT_TIME format: %v", err)
	}

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
		DailyReportTime:        dailyReportTime,
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

// parseTime parses a time string in HH:MM format and returns a time.Duration
func parseTime(timeStr string) (time.Duration, error) {
	var hour, minute int
	_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	if err != nil {
		return 0, err
	}

	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	return time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute, nil
}

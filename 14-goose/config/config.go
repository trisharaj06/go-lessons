package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBConnString    string
	Port            string
	AppBaseURL      string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	MagicLinkTTL    time.Duration
	CookieSecure    bool
	CookieDomain    string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func Load() (*Config, error) {
	dbConnString := os.Getenv("GOOSE_DBSTRING")
	if dbConnString == "" {
		return nil, fmt.Errorf("GOOSE_DBSTRING is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	cfg := &Config{
		DBConnString:    dbConnString,
		Port:            envOrDefault("APP_PORT", "8080"),
		AppBaseURL:      envOrDefault("APP_BASE_URL", "http://localhost:8080"),
		JWTSecret:       jwtSecret,
		AccessTokenTTL:  minutes(envOrDefault("ACCESS_TOKEN_TTL_MINUTES", "15")),
		RefreshTokenTTL: hours(envOrDefault("REFRESH_TOKEN_TTL_HOURS", "168")),
		MagicLinkTTL:    minutes(envOrDefault("MAGIC_LINK_TTL_MINUTES", "15")),
		CookieSecure:    envOrDefault("COOKIE_SECURE", "false") == "true",
		CookieDomain:    os.Getenv("COOKIE_DOMAIN"),

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     envOrDefault("SMTP_PORT", "587"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func minutes(raw string) time.Duration {
	n, err := strconv.Atoi(raw)
	if err != nil {
		n = 15
	}
	return time.Duration(n) * time.Minute
}

func hours(raw string) time.Duration {
	n, err := strconv.Atoi(raw)
	if err != nil {
		n = 168
	}
	return time.Duration(n) * time.Hour
}

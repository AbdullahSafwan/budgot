package configs

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port        string
	AppEnv      string
	CSRFAuthKey string
	LogLevel    string
	DatabaseDSN string

	SessionTTL           time.Duration
	SessionIdleTimeout   time.Duration
	BcryptCost           int
	LoginRateLimit       int
	LoginRateLimitWindow time.Duration
	TrustedProxyHops     int
}

func (c Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func requireEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("%s is not set", key)
	}
	return v, nil
}

func optionalEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func LoadConfig() (Config, error) {
	csrfKey, err := requireEnv("CSRF_AUTH_KEY")
	if err != nil {
		return Config{}, err
	}
	if len(csrfKey) < 32 {
		return Config{}, fmt.Errorf("CSRF_AUTH_KEY must be at least 32 bytes, got %d", len(csrfKey))
	}
	port, err := requireEnv("PORT")
	if err != nil {
		return Config{}, err
	}
	appEnv, err := requireEnv("APP_ENV")
	if err != nil {
		return Config{}, err
	}
	logLevel, err := requireEnv("LOG_LEVEL")
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:        port,
		AppEnv:      appEnv,
		CSRFAuthKey: csrfKey,
		LogLevel:    logLevel,
		DatabaseDSN: optionalEnv("DATABASE_DSN", "file:budgot.db"),

		SessionTTL:           7 * 24 * time.Hour,
		SessionIdleTimeout:   30 * time.Minute,
		BcryptCost:           12,
		LoginRateLimit:       5,
		LoginRateLimitWindow: time.Minute,
		TrustedProxyHops:     0,
	}, nil
}

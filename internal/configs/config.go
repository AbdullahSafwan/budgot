package configs

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	AppEnv      string
	CSRFAuthKey string
	LogLevel    string
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

func LoadConfig() (Config, error) {
	csrfKey, err := requireEnv("CSRF_AUTH_KEY")
	if err != nil {
		return Config{}, err
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
	}, nil
}

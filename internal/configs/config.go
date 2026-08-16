package configs

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	AppEnv      string
	CSRFAuthKey string
}

func (c Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func LoadConfig() (Config, error) {
	csrfKey := os.Getenv("CSRF_AUTH_KEY")
	if csrfKey == "" {
		return Config{}, fmt.Errorf("CSRF_AUTH_KEY is not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	return Config{
		Port:        port,
		AppEnv:      appEnv,
		CSRFAuthKey: csrfKey,
	}, nil
}

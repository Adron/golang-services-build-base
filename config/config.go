package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            int
	DatadogAgentURL string
	LogLevel        string
}

func LoadConfig() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))

	return &Config{
		Port:            port,
		DatadogAgentURL: getEnv("DD_AGENT_URL", "127.0.0.1:8125"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

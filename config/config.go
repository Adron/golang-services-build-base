package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port             int
	OtelEndpoint     string
	LogLevel         string
	ServiceName      string
	ServiceVersion   string
	ServiceNamespace string
}

func LoadConfig() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))

	return &Config{
		Port:             port,
		OtelEndpoint:     getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		ServiceName:      getEnv("SERVICE_NAME", "vision-service"),
		ServiceVersion:   getEnv("SERVICE_VERSION", "1.0.0"),
		ServiceNamespace: getEnv("SERVICE_NAMESPACE", "default"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

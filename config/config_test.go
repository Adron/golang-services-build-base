package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expected: &Config{
				Port:             8080,
				OtelEndpoint:     "http://localhost:4318",
				LogLevel:         "info",
				ServiceName:      "vision-service",
				ServiceVersion:   "1.0.0",
				ServiceNamespace: "default",
			},
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"PORT":              "9090",
				"LOG_LEVEL":         "debug",
				"SERVICE_NAME":      "test-service",
				"SERVICE_VERSION":   "2.0.0",
				"SERVICE_NAMESPACE": "test",
			},
			expected: &Config{
				Port:             9090,
				OtelEndpoint:     "http://localhost:4318",
				LogLevel:         "debug",
				ServiceName:      "test-service",
				ServiceVersion:   "2.0.0",
				ServiceNamespace: "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Clearenv()

			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Load configuration
			cfg := LoadConfig()

			// Compare values
			if cfg.Port != tt.expected.Port {
				t.Errorf("Port = %v, want %v", cfg.Port, tt.expected.Port)
			}
			if cfg.OtelEndpoint != tt.expected.OtelEndpoint {
				t.Errorf("OtelEndpoint = %v, want %v", cfg.OtelEndpoint, tt.expected.OtelEndpoint)
			}
			if cfg.LogLevel != tt.expected.LogLevel {
				t.Errorf("LogLevel = %v, want %v", cfg.LogLevel, tt.expected.LogLevel)
			}
			if cfg.ServiceName != tt.expected.ServiceName {
				t.Errorf("ServiceName = %v, want %v", cfg.ServiceName, tt.expected.ServiceName)
			}
			if cfg.ServiceVersion != tt.expected.ServiceVersion {
				t.Errorf("ServiceVersion = %v, want %v", cfg.ServiceVersion, tt.expected.ServiceVersion)
			}
			if cfg.ServiceNamespace != tt.expected.ServiceNamespace {
				t.Errorf("ServiceNamespace = %v, want %v", cfg.ServiceNamespace, tt.expected.ServiceNamespace)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		defaultVal string
		envValue   string
		expected   string
	}{
		{
			name:       "existing environment variable",
			key:        "TEST_VAR",
			defaultVal: "default",
			envValue:   "test_value",
			expected:   "test_value",
		},
		{
			name:       "non-existing environment variable",
			key:        "NON_EXISTING_VAR",
			defaultVal: "default",
			envValue:   "",
			expected:   "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Clearenv()

			// Set test environment variable if needed
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			// Test getEnv function
			result := getEnv(tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("getEnv(%v, %v) = %v, want %v", tt.key, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

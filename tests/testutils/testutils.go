package testutils

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

// TestContext returns a context with a test tracer
func TestContext(t *testing.T) context.Context {
	ctx := context.Background()
	tracer := otel.Tracer("test")
	ctx, _ = tracer.Start(ctx, "test-span")
	return ctx
}

// TestServer creates a test HTTP server with the given handler
func TestServer(t *testing.T, handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

// AssertResponse checks if the response matches the expected status code and body
func AssertResponse(t *testing.T, resp *http.Response, expectedStatus int, expectedBody string) {
	assert.Equal(t, expectedStatus, resp.StatusCode)
	if expectedBody != "" {
		body := readBody(t, resp)
		assert.Equal(t, expectedBody, body)
	}
}

// LoadTestConfig returns a configuration suitable for load testing
func LoadTestConfig() map[string]string {
	return map[string]string{
		"PORT":            "8080",
		"OTEL_ENDPOINT":   "http://localhost:4318",
		"LOG_LEVEL":       "error", // Reduce logging during load tests
		"SERVICE_NAME":    "vision-service-test",
		"SERVICE_VERSION": "test",
	}
}

// BenchmarkConfig returns a configuration suitable for benchmarks
func BenchmarkConfig() map[string]string {
	return map[string]string{
		"PORT":            "8080",
		"OTEL_ENDPOINT":   "http://localhost:4318",
		"LOG_LEVEL":       "error",
		"SERVICE_NAME":    "vision-service-bench",
		"SERVICE_VERSION": "bench",
	}
}

// WaitForServer waits for the server to be ready
func WaitForServer(t *testing.T, url string, timeout time.Duration) {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			t.Fatal("server did not become ready in time")
		}
		resp, err := http.Get(url + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// readBody reads and returns the response body
func readBody(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	return string(body)
}

// ReadBody reads the response body and returns it as a string
func ReadBody(t *testing.T, resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	return string(body)
}

// StartTestServer starts a test server with the given configuration
func StartTestServer(t *testing.T, config map[string]string) *httptest.Server {
	// Set environment variables from config
	for k, v := range config {
		t.Setenv(k, v)
	}

	// Create a new router
	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// Create and start the server
	server := httptest.NewServer(router)
	return server
}

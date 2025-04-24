package load

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adron/golang-services-build-base/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandlerLoad(t *testing.T) {
	handler := handlers.NewHealthHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	// Test parameters
	concurrency := 100
	duration := 10 * time.Second
	errors := make(chan error, concurrency)

	// Start load test
	start := time.Now()
	for i := 0; i < concurrency; i++ {
		go func() {
			for time.Since(start) < duration {
				resp, err := http.Get(server.URL + "/health")
				if err != nil {
					errors <- err
					continue
				}
				if resp.StatusCode != http.StatusOK {
					errors <- err
				}
				resp.Body.Close()
			}
		}()
	}

	// Wait for test duration
	time.Sleep(duration)

	// Check for errors
	select {
	case err := <-errors:
		t.Fatalf("Error during load test: %v", err)
	default:
		// No errors, test passed
	}

	// Verify server is still responding
	resp, err := http.Get(server.URL + "/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

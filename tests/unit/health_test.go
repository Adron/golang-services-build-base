package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adron/golang-services-build-base/internal/handlers"
	"github.com/adron/golang-services-build-base/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"healthy","timestamp":`,
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			handler := handlers.NewHealthHandler()
			handler.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			if tt.expectedBody != "" {
				body := testutils.ReadBody(t, resp)
				assert.Contains(t, body, tt.expectedBody)
			}
		})
	}
}

func TestHealthHandlerConcurrent(t *testing.T) {
	handler := handlers.NewHealthHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	concurrency := 100
	requests := 1000
	errors := make(chan error, concurrency*requests)

	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < requests; j++ {
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

	// Wait for all requests to complete
	time.Sleep(2 * time.Second)

	select {
	case err := <-errors:
		t.Fatalf("Error during concurrent requests: %v", err)
	default:
		// No errors, test passed
	}
}

func TestHealthHandlerTimeout(t *testing.T) {
	// Create a slow handler that takes longer than the client timeout
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Simulate slow processing
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create a client with a short timeout
	client := &http.Client{
		Timeout: 1 * time.Millisecond,
	}

	// Make the request
	_, err := client.Get(server.URL)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

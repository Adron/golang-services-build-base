package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTestContext(t *testing.T) {
	ctx := TestContext(t)
	assert.NotNil(t, ctx)
}

func TestTestServer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := TestServer(t, handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAssertResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.NoError(t, err)
	AssertResponse(t, resp, http.StatusOK, "test response")
}

func TestBenchmarkConfig(t *testing.T) {
	config := BenchmarkConfig()
	assert.NotEmpty(t, config)
	assert.Equal(t, "8080", config["PORT"])
	assert.Equal(t, "http://localhost:4318", config["OTEL_ENDPOINT"])
}

func TestReadBody(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test body"))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.NoError(t, err)
	body := ReadBody(t, resp)
	assert.Equal(t, "test body", body)
}

func TestWaitForServer(t *testing.T) {
	// Test successful case
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Give the server more time to start
	time.Sleep(500 * time.Millisecond)

	// Try multiple times with increasing timeouts
	for i := 0; i < 3; i++ {
		timeout := time.Duration(i+1) * time.Second
		err := WaitForServer(t, server.URL, timeout)
		if err == nil {
			break
		}
		if i == 2 {
			t.Fatalf("Server did not become ready after multiple attempts")
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Test timeout case with a non-existent server
	start := time.Now()
	WaitForServer(t, "http://localhost:9999", 100*time.Millisecond)
	elapsed := time.Since(start)
	assert.GreaterOrEqual(t, elapsed, 100*time.Millisecond, "Should wait for at least the timeout duration")
}

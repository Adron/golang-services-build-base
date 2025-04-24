package load

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/adron/golang-services-build-base/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestLoadHealthEndpoint(t *testing.T) {
	config := testutils.LoadTestConfig()
	server := testutils.StartTestServer(t, config)
	defer server.Close()

	testutils.WaitForServer(t, server.URL, 5*time.Second)

	// Test parameters
	concurrency := 50
	duration := 30 * time.Second
	requests := make(chan int, concurrency)
	errors := make(chan error, concurrency)

	// Start workers
	for i := 0; i < concurrency; i++ {
		go func() {
			client := &http.Client{
				Timeout: 5 * time.Second,
			}
			for range requests {
				resp, err := client.Get(server.URL + "/health")
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

	// Send requests for the specified duration
	start := time.Now()
	requestCount := 0
	for time.Since(start) < duration {
		select {
		case requests <- requestCount:
			requestCount++
		case err := <-errors:
			t.Fatalf("Error during load test: %v", err)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	close(requests)

	// Calculate metrics
	elapsed := time.Since(start)
	requestsPerSecond := float64(requestCount) / elapsed.Seconds()

	t.Logf("Load test results:")
	t.Logf("Duration: %v", elapsed)
	t.Logf("Total requests: %d", requestCount)
	t.Logf("Requests per second: %.2f", requestsPerSecond)

	// Assert minimum performance requirements
	assert.Greater(t, requestsPerSecond, 100.0, "Should handle at least 100 requests per second")
	assert.Less(t, elapsed, 35*time.Second, "Test should complete within 35 seconds")
}

func TestLoadWithIncreasingConcurrency(t *testing.T) {
	config := testutils.LoadTestConfig()
	server := testutils.StartTestServer(t, config)
	defer server.Close()

	testutils.WaitForServer(t, server.URL, 5*time.Second)

	concurrencyLevels := []int{10, 25, 50, 100}
	duration := 10 * time.Second

	for _, concurrency := range concurrencyLevels {
		t.Run(fmt.Sprintf("concurrency-%d", concurrency), func(t *testing.T) {
			requests := make(chan int, concurrency)
			errors := make(chan error, concurrency)

			// Start workers
			for i := 0; i < concurrency; i++ {
				go func() {
					client := &http.Client{
						Timeout: 5 * time.Second,
					}
					for range requests {
						resp, err := client.Get(server.URL + "/health")
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

			// Send requests for the specified duration
			start := time.Now()
			requestCount := 0
			for time.Since(start) < duration {
				select {
				case requests <- requestCount:
					requestCount++
				case err := <-errors:
					t.Fatalf("Error during load test: %v", err)
				default:
					time.Sleep(10 * time.Millisecond)
				}
			}
			close(requests)

			// Calculate metrics
			elapsed := time.Since(start)
			requestsPerSecond := float64(requestCount) / elapsed.Seconds()

			t.Logf("Load test results for concurrency %d:", concurrency)
			t.Logf("Duration: %v", elapsed)
			t.Logf("Total requests: %d", requestCount)
			t.Logf("Requests per second: %.2f", requestsPerSecond)

			// Assert minimum performance requirements
			assert.Greater(t, requestsPerSecond, float64(concurrency),
				"Should handle at least %d requests per second", concurrency)
		})
	}
}

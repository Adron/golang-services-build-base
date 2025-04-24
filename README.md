# Computer Vision Service

A Windows-based service for computer vision capabilities to identify lines of people and vehicles for order processing.

## Features

- Enterprise-grade observability with OpenTelemetry integration
- Structured logging
- Performance metrics
- Health check endpoints
- Configurable through environment variables
- Comprehensive test suite with performance benchmarks
- Load testing capabilities
- Test coverage reporting and monitoring

## Prerequisites

- Go 1.23 or later
- Docker and Docker Compose (for OpenTelemetry Collector)
- Windows OS (for service deployment)

## Configuration

The service can be configured using the following environment variables:

- `PORT`: Service port (default: 8080)
- `OTEL_ENDPOINT`: OpenTelemetry Collector endpoint (default: http://localhost:4318)
- `LOG_LEVEL`: Logging level (default: info)
- `SERVICE_NAME`: Service name for telemetry (default: vision-service)
- `SERVICE_VERSION`: Service version (default: 1.0.0)
- `SERVICE_NAMESPACE`: Service namespace (default: default)

## Running the Service

1. Start the OpenTelemetry Collector:
   ```bash
   docker-compose up -d
   ```

2. Set any required environment variables
3. Run the service:
   ```bash
   go run main.go
   ```

## Health Check

The service exposes a health check endpoint at `/health` that returns a 200 OK status when the service is running properly.

## Observability

The service uses OpenTelemetry for comprehensive observability:

### Traces
- Request tracing
- Service operation spans
- Distributed tracing support

### Metrics
- Service uptime
- Request counts
- Response times
- Error rates
- Custom business metrics

### Logging
Logs are structured in JSON format, including:
- Service startup/shutdown events
- Request/response information
- Error details
- Performance metrics

## Testing

The service includes a comprehensive test suite with multiple testing layers:

### Test Structure
```
tests/
├── unit/           # Unit tests for individual components
├── integration/    # Integration tests for service interactions
├── benchmarks/     # Performance benchmarks
└── load/          # Load testing scenarios
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
./scripts/coverage.sh

# Run benchmarks
go test -bench=. -benchmem ./...

# Run load tests
go test -v ./tests/load/...
```

### Test Coverage

The project maintains a minimum test coverage threshold of 80%. Coverage reports are generated using the `coverage.sh` script:

```bash
./scripts/coverage.sh
```

This will generate:
- HTML coverage report (`coverage/coverage.html`)
- Text summary (`coverage/coverage.txt`)
- Coverage badge (`coverage/badge.svg`)

### Load Testing

The service includes comprehensive load testing capabilities:

1. **Basic Load Test**
   - 50 concurrent users
   - 30-second duration
   - Measures requests per second
   - Validates response times

2. **Concurrency Scaling Test**
   - Tests with 10, 25, 50, and 100 concurrent users
   - 10-second duration per level
   - Measures performance degradation
   - Identifies bottlenecks

### Benchmark Results

The service includes performance benchmarks that measure:
- Health check endpoint performance
- Concurrent request handling
- Memory allocation patterns

Example benchmark output:
```
BenchmarkHealthCheck-10          748,286 ops/s    1457 ns/op    1413 B/op
BenchmarkHealthCheckParallel-10  1,396,557 ops/s  786.1 ns/op   1428 B/op
```

These benchmarks help ensure the service maintains optimal performance under load.

## Testing Infrastructure

The service includes a comprehensive testing infrastructure with multiple layers of testing to ensure reliability and performance.

### Test Structure

```
tests/
├── unit/           # Unit tests for individual components
├── integration/    # Integration tests for service interactions
├── benchmarks/     # Performance benchmarks
└── load/          # Load testing scenarios
```

### Unit Tests

Unit tests focus on testing individual components in isolation:

```bash
# Run all unit tests
go test -v ./tests/unit/...

# Run specific unit test
go test -v ./tests/unit/health_test.go
```

Key unit test features:
- Isolated component testing
- Mock dependencies
- Concurrent request handling
- Error scenarios
- Timeout handling

Example unit test structure:
```go
func TestHealthHandler(t *testing.T) {
    tests := []struct {
        name           string
        method         string
        expectedStatus int
        expectedBody   string
    }{
        // Test cases
    }
    // Test execution
}
```

### Integration Tests

Integration tests verify component interactions:

```bash
# Run all integration tests
go test -v ./tests/integration/...
```

Integration test features:
- Component interaction testing
- End-to-end request flow
- Database interactions
- External service mocking
- Configuration validation

### Load Testing

Load tests verify system performance under stress:

```bash
# Run all load tests
go test -v ./tests/load/...
```

Load test scenarios:
1. **Basic Load Test**
   - 50 concurrent users
   - 30-second duration
   - Measures requests per second
   - Validates response times

2. **Concurrency Scaling Test**
   - Tests with 10, 25, 50, and 100 concurrent users
   - 10-second duration per level
   - Measures performance degradation
   - Identifies bottlenecks

Load test metrics:
- Requests per second
- Response times
- Error rates
- Resource utilization
- Memory allocation

### Code Coverage

The project maintains a minimum test coverage threshold of 80%. Coverage reports are generated using the `coverage.sh` script:

```bash
# Generate coverage reports
./scripts/coverage.sh
```

This generates:
- HTML coverage report (`coverage/coverage.html`)
- Text summary (`coverage/coverage.txt`)
- Coverage badge (`coverage/badge.svg`)

Coverage features:
- Line-by-line coverage analysis
- Function coverage metrics
- Package-level coverage
- Coverage thresholds
- Visual coverage reports

### Test Utilities

The project includes a test utilities package (`tests/testutils/`) with common testing functions:

```go
// Test context creation
func TestContext(t *testing.T) context.Context

// Test server setup
func TestServer(t *testing.T, handler http.Handler) *httptest.Server

// Response assertions
func AssertResponse(t *testing.T, resp *http.Response, expectedStatus int, expectedBody string)

// Load test configuration
func LoadTestConfig() map[string]string

// Benchmark configuration
func BenchmarkConfig() map[string]string

// Server readiness check
func WaitForServer(t *testing.T, url string, timeout time.Duration)
```

### Running All Tests

To run the complete test suite:

```bash
# Run all tests with coverage
./scripts/coverage.sh

# Run specific test types
go test -v ./tests/unit/...    # Unit tests
go test -v ./tests/integration/...  # Integration tests
go test -v ./tests/load/...    # Load tests
go test -bench=. -benchmem ./...  # Benchmarks
```

### Test Configuration

Tests can be configured using environment variables:

```bash
# Test configuration
export TEST_TIMEOUT=30s
export TEST_CONCURRENCY=50
export TEST_DURATION=30s

# Coverage configuration
export MIN_COVERAGE=80
export COVERAGE_MODE=atomic
```

### Continuous Integration

The test suite is integrated into the CI pipeline:
- Automated test execution
- Coverage threshold enforcement
- Performance regression detection
- Load test results analysis 
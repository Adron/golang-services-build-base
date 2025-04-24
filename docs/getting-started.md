# Getting Started with the Computer Vision Service

This guide will help you set up and run the Computer Vision Service locally for development and testing.

## Prerequisites

Before you begin, ensure you have the following installed:

- Go 1.23 or later
- Docker and Docker Compose
- Windows OS (for service deployment)
- Git

## Project Structure

```
vision-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── metrics/
│   └── tracing/
├── pkg/
│   └── health/
├── docs/
├── tests/
├── docker-compose.yml
└── README.md
```

## Setting Up the Project

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/vision-service.git
   cd vision-service
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Start the OpenTelemetry Collector:
   ```bash
   docker-compose up -d
   ```

## Running the Service

### Interactive Mode

1. Set environment variables:
   ```bash
   export PORT=8080
   export OTEL_ENDPOINT=http://localhost:4318
   export LOG_LEVEL=debug
   ```

2. Run the service:
   ```bash
   go run cmd/server/main.go
   ```

### Headless Mode

1. Build the service:
   ```bash
   go build -o vision-service cmd/server/main.go
   ```

2. Run the service:
   ```bash
   ./vision-service
   ```

## Configuration

The service can be configured using environment variables:

- `PORT`: Service port (default: 8080)
- `OTEL_ENDPOINT`: OpenTelemetry Collector endpoint (default: http://localhost:4318)
- `LOG_LEVEL`: Logging level (default: info)
- `SERVICE_NAME`: Service name for telemetry (default: vision-service)
- `SERVICE_VERSION`: Service version (default: 1.0.0)
- `SERVICE_NAMESPACE`: Service namespace (default: default)

## Health Check

The service exposes a health check endpoint at `/health`. You can test it using curl:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2024-03-21T12:00:00Z"
}
```

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

## Development Workflow

1. Make changes to the code
2. Run tests:
   ```bash
   go test -v ./...
   ```
3. Run benchmarks:
   ```bash
   go test -bench=. -benchmem ./...
   ```
4. Build and test the service
5. Commit changes and create a pull request

## Troubleshooting

### Common Issues

1. **Service fails to start**
   - Check if the OpenTelemetry Collector is running
   - Verify environment variables are set correctly
   - Check port availability

2. **No metrics or logs**
   - Ensure OpenTelemetry Collector is running
   - Verify the collector endpoint is correct
   - Check network connectivity

3. **High latency**
   - Check system resources
   - Review benchmark results
   - Monitor OpenTelemetry metrics

## Next Steps

1. Review the API documentation
2. Set up monitoring dashboards
3. Configure alerts based on metrics
4. Set up CI/CD pipeline
5. Deploy to production environment 
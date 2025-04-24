# Computer Vision Service

A Windows-based service for computer vision capabilities to identify lines of people and vehicles for order processing.

## Features

- Enterprise-grade observability with Datadog integration
- Structured logging
- Performance metrics
- Health check endpoints
- Configurable through environment variables

## Prerequisites

- Go 1.23 or later
- Datadog Agent running locally (for metrics and logging)
- Windows OS (for service deployment)

## Configuration

The service can be configured using the following environment variables:

- `PORT`: Service port (default: 8080)
- `DD_AGENT_URL`: Datadog Agent URL (default: 127.0.0.1:8125)
- `LOG_LEVEL`: Logging level (default: info)

## Running the Service

1. Ensure the Datadog Agent is running
2. Set any required environment variables
3. Run the service:
   ```bash
   go run main.go
   ```

## Health Check

The service exposes a health check endpoint at `/health` that returns a 200 OK status when the service is running properly.

## Metrics

The service sends the following metrics to Datadog:
- Service uptime
- Request counts
- Response times
- Error rates

## Logging

Logs are sent to Datadog in JSON format, including:
- Service startup/shutdown events
- Request/response information
- Error details
- Performance metrics 
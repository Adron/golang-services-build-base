# Getting Started with the Computer Vision Service Base

*Published: April 23, 2025*

Hey there! I'm Adron Hall, and today I'm going to walk you through setting up this computer vision service base project. This service is designed to run as a Windows service, providing computer vision capabilities for identifying lines of people and vehicles in order processing scenarios. Let's dive in and get this thing set up!

## Repository

You can find this project at:
```
git@github.com:Adron/golang-services-build-base.git
```

## Prerequisites

Before we get started, make sure you've got these tools installed:

- Go 1.23 or later (we're using the latest and greatest!)
- A Windows environment (since this is designed to run as a Windows service)
- Datadog Agent (for metrics and logging)

## Project Structure

Let's take a quick look at what we've built:

```
golang-services-build-base/
├── main.go              # Main service entry point
├── config/              # Configuration management
│   └── config.go        # Configuration structs and loading
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
└── .gitignore          # Git ignore patterns
```

## Setting Up the Project

### 1. Clone the Repository

First, let's get the code:

```bash
git clone git@github.com:Adron/golang-services-build-base.git
cd golang-services-build-base
```

### 2. Initialize the Go Module

The module is already set up with the correct path:

```bash
go mod init github.com/Adron/golang-services-build-base
```

### 3. Add Dependencies

We're using some solid libraries for this project:

- `github.com/DataDog/datadog-go/v5` - For metrics and logging
- `github.com/gorilla/mux` - For HTTP routing
- `github.com/sirupsen/logrus` - For structured logging
- `github.com/spf13/cobra` - For CLI interface

Run this to get everything downloaded:

```bash
go mod tidy
```

### 4. Build the Service

Let's build our service:

```bash
go build -o vision-service
```

## Running the Service

You've got two ways to run this service:

### Interactive Mode

This is the default mode, giving you a nice terminal interface to control the service:

```bash
./vision-service
```

You'll see a menu with these options:
- Press 's' to start the service
- Press 'q' to stop the service
- Press 'x' to exit

### Headless Mode

For production or when you want to run it without the interactive interface:

```bash
./vision-service --headless
```

## Configuration

The service is configurable through environment variables:

- `PORT`: Service port (default: 8080)
- `DD_AGENT_URL`: Datadog Agent URL (default: 127.0.0.1:8125)
- `LOG_LEVEL`: Logging level (default: info)

## Health Check

Once the service is running, you can check its health at:

```
http://localhost:8080/health
```

## Observability

We've baked in some serious observability features:

- Structured logging with Logrus
- Metrics collection with Datadog
- Health check endpoints
- Performance monitoring

## Development Workflow

Here's my typical workflow when working on this service:

1. Make changes to the code
2. Run tests (when we add them)
3. Build the service: `go build -o vision-service`
4. Test in interactive mode: `./vision-service`
5. Deploy to test environment in headless mode: `./vision-service --headless`

## Next Steps

Looking to extend this base? Here are some ideas:

1. Add computer vision integration (OpenCV, TensorFlow, etc.)
2. Implement Windows service installation
3. Add more metrics and logging
4. Create a proper CI/CD pipeline

## Troubleshooting

If you run into issues:

1. Check the logs - they're in JSON format for easy parsing
2. Verify the Datadog Agent is running
3. Make sure no other service is using port 8080
4. Check your environment variables

## Wrapping Up

That's it! You've now got a solid foundation for a computer vision service with enterprise-grade observability. The service is ready to be extended with actual computer vision capabilities while maintaining good operational practices.

Got questions? Hit me up on [Twitter](https://twitter.com/adron) or check out more of my writing at [Composite Code](https://compositecode.blog/).

Happy coding! 🚀 
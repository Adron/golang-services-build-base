package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/adron/golang-services-build-base/config"
)

var (
	logger *logrus.Logger
	server *http.Server
	cfg    *config.Config
	meter  otelmetric.Meter
	tracer trace.Tracer
)

func init() {
	// Load configuration
	cfg = config.LoadConfig()

	// Initialize logger
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	// Initialize OpenTelemetry
	ctx := context.Background()

	// Create OTLP trace exporter
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		logger.Fatalf("Failed to create trace exporter: %v", err)
	}

	// Create OTLP metric exporter
	metricExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint("localhost:4318"),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		logger.Fatalf("Failed to create metric exporter: %v", err)
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.ServiceNamespace(cfg.ServiceNamespace),
		),
	)
	if err != nil {
		logger.Fatalf("Failed to create resource: %v", err)
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer(cfg.ServiceName)

	// Create meter provider
	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(mp)
	meter = mp.Meter(cfg.ServiceName)
}

func startServer() {
	// Create a new router
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "health-check")
		defer span.End()

		// Record metrics
		counter, err := meter.Int64Counter("health_check_count")
		if err != nil {
			logger.Errorf("Failed to create counter: %v", err)
		} else {
			counter.Add(ctx, 1)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Service is healthy")
	}).Methods("GET")

	// Create HTTP server
	server = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting service on port %d", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()
}

func stopServer() {
	if server != nil {
		logger.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Server forced to shutdown: %v", err)
		}
		logger.Info("Server stopped")
	}
}

func runService(headless bool) {
	if headless {
		// Start server in headless mode
		startServer()
		logger.Info("Service running in headless mode")

		// Wait for interrupt signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		stopServer()
		return
	}

	// Interactive mode
	fmt.Print("\033[H\033[2J")
	fmt.Println("Computer Vision Service Control")
	fmt.Println("==============================")
	fmt.Println("Press 's' to start the service")
	fmt.Println("Press 'q' to stop the service")
	fmt.Println("Press 'x' to exit")
	fmt.Println("==============================")

	// Channel to receive keyboard input
	inputChan := make(chan string)
	go func() {
		var input string
		for {
			fmt.Scanln(&input)
			inputChan <- input
		}
	}()

	// Main control loop
	for {
		select {
		case cmd := <-inputChan:
			switch cmd {
			case "s":
				if server == nil {
					startServer()
					fmt.Println("Service started successfully")
				} else {
					fmt.Println("Service is already running")
				}
			case "q":
				if server != nil {
					stopServer()
					server = nil
					fmt.Println("Service stopped successfully")
				} else {
					fmt.Println("Service is not running")
				}
			case "x":
				if server != nil {
					stopServer()
				}
				fmt.Println("Exiting...")
				return
			default:
				fmt.Println("Invalid command. Use 's' to start, 'q' to stop, or 'x' to exit")
			}
		}
	}
}

func main() {
	var headless bool

	rootCmd := &cobra.Command{
		Use:   "vision-service",
		Short: "Computer Vision Service for line detection",
		Long: `A Windows-based service for computer vision capabilities to identify lines of people 
and vehicles for order processing.`,
		Run: func(cmd *cobra.Command, args []string) {
			runService(headless)
		},
	}

	rootCmd.Flags().BoolVarP(&headless, "headless", "H", false, "Run service in headless mode")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

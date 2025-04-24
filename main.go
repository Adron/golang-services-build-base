package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logger *logrus.Logger
	stats  *statsd.Client
	server *http.Server
)

func init() {
	// Initialize logger
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	// Initialize Datadog statsd client
	var err error
	stats, err = statsd.New("127.0.0.1:8125")
	if err != nil {
		logger.Fatalf("Failed to create statsd client: %v", err)
	}
}

func startServer() {
	// Create a new router
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Service is healthy")
	}).Methods("GET")

	// Create HTTP server
	server = &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting service on :8080")
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

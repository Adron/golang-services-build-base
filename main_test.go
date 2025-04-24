package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func setupTestTracer() (*trace.TracerProvider, error) {
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("test")
	return tp, nil
}

func setupTestMeter() (*metric.MeterProvider, error) {
	mp := metric.NewMeterProvider()
	otel.SetMeterProvider(mp)
	meter = mp.Meter("test")
	return mp, nil
}

func TestHealthCheckHandler(t *testing.T) {
	// Setup test tracer and meter
	tp, err := setupTestTracer()
	if err != nil {
		t.Fatalf("Failed to setup test tracer: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			t.Errorf("Failed to shutdown test tracer: %v", err)
		}
	}()

	mp, err := setupTestMeter()
	if err != nil {
		t.Fatalf("Failed to setup test meter: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := mp.Shutdown(ctx); err != nil {
			t.Errorf("Failed to shutdown test meter: %v", err)
		}
	}()

	// Create test router and register handler
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "health-check")
		defer span.End()

		counter, err := meter.Int64Counter("health_check_count")
		if err != nil {
			t.Errorf("Failed to create counter: %v", err)
		} else {
			counter.Add(ctx, 1)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is healthy"))
	}).Methods("GET")

	// Create test request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Health check returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	expected := "Service is healthy"
	if w.Body.String() != expected {
		t.Errorf("Health check returned unexpected body: got %v want %v", w.Body.String(), expected)
	}
}

func BenchmarkHealthCheck(b *testing.B) {
	// Setup test tracer and meter
	tp, _ := setupTestTracer()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		tp.Shutdown(ctx)
	}()

	mp, _ := setupTestMeter()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		mp.Shutdown(ctx)
	}()

	// Create test router and register handler
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "health-check")
		defer span.End()

		counter, _ := meter.Int64Counter("health_check_count")
		counter.Add(ctx, 1)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is healthy"))
	}).Methods("GET")

	// Create test request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkHealthCheckParallel(b *testing.B) {
	// Setup test tracer and meter
	tp, _ := setupTestTracer()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		tp.Shutdown(ctx)
	}()

	mp, _ := setupTestMeter()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		mp.Shutdown(ctx)
	}()

	// Create test router and register handler
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "health-check")
		defer span.End()

		counter, _ := meter.Int64Counter("health_check_count")
		counter.Add(ctx, 1)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is healthy"))
	}).Methods("GET")

	// Run parallel benchmark
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		for pb.Next() {
			router.ServeHTTP(w, req)
		}
	})
}

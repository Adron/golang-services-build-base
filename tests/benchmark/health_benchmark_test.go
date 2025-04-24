package benchmark

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adron/golang-services-build-base/internal/handlers"
)

func BenchmarkHealthHandler(b *testing.B) {
	handler := handlers.NewHealthHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(w, req)
	}
}

func BenchmarkHealthHandlerConcurrent(b *testing.B) {
	handler := handlers.NewHealthHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Get(server.URL + "/health")
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})
}

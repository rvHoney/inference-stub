package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestChatCompletionsHandler(t *testing.T) {
	cfg := HandlerConfig{
		TTFT: 10 * time.Millisecond,
		TPOT: 5 * time.Millisecond,
	}

	req, err := http.NewRequest("POST", "/v1/chat/completions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := ChatCompletionsHandler(cfg)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check response headers for SSE
	expectedHeaders := map[string]string{
		"Content-Type":  "text/event-stream",
		"Cache-Control": "no-cache",
		"Connection":    "keep-alive",
	}

	for key, expectedValue := range expectedHeaders {
		if value := rr.Header().Get(key); value != expectedValue {
			t.Errorf("handler returned wrong header for %s: got %v want %v",
				key, value, expectedValue)
		}
	}
}

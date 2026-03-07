package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/robin-vidal/inference-stub/pkg/lorem"
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
		TTFT:      10 * time.Millisecond,
		TPOT:      5 * time.Millisecond,
		Generator: lorem.New(50),
	}

	reqBody := `{"model": "test-model", "messages": [], "stream": true}`
	req, err := http.NewRequest("POST", "/v1/chat/completions", strings.NewReader(reqBody))
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

func TestChatCompletionsHandler_NoStream(t *testing.T) {
	cfg := HandlerConfig{
		TTFT:      5 * time.Millisecond,
		TPOT:      2 * time.Millisecond,
		Generator: lorem.New(10),
	}

	reqBody := `{"model": "test-model", "messages": [{"role":"user", "content":"hello"}], "stream": false}`
	req, err := http.NewRequest("POST", "/v1/chat/completions", strings.NewReader(reqBody))
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

	expectedContentType := "application/json"
	if value := rr.Header().Get("Content-Type"); value != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", value, expectedContentType)
	}

	var resp ChatCompletionResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Model != "test-model" {
		t.Errorf("expected model %v, got %v", "test-model", resp.Model)
	}
	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %v", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content == "" {
		t.Errorf("unexpected empty content")
	}
	if resp.Usage.CompletionTokens != 10 {
		t.Errorf("expected 10 completion tokens, got %v", resp.Usage.CompletionTokens)
	}
}

func TestChatCompletionsHandler_InvalidJSON(t *testing.T) {
	cfg := HandlerConfig{
		Generator: lorem.New(50),
	}

	reqBody := `{"model": "test-model", "messages": ` // malformed JSON
	req, err := http.NewRequest("POST", "/v1/chat/completions", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := ChatCompletionsHandler(cfg)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestChatCompletionsHandler_StreamOutput(t *testing.T) {
	cfg := HandlerConfig{
		TTFT:      2 * time.Millisecond,
		TPOT:      1 * time.Millisecond,
		Generator: lorem.New(50),
	}

	reqBody := `{"model": "test-model", "messages": [], "stream": true}`
	req, err := http.NewRequest("POST", "/v1/chat/completions", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := ChatCompletionsHandler(cfg)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	bodyStr := rr.Body.String()

	if !strings.Contains(bodyStr, "data: {") {
		t.Errorf("Stream response does not contain JSON 'data:' chunks")
	}

	if !strings.Contains(bodyStr, "data: [DONE]\n\n") {
		t.Errorf("Stream response is missing the final [DONE] signal")
	}
}

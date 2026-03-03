package server

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestServerLifecycle(t *testing.T) {
	port := 18080 // Different port to avoid conflicts
	timeout := 1 * time.Second
	ttft := 10 * time.Millisecond
	tpot := 10 * time.Millisecond

	srv := Init(port, timeout, ttft, tpot)

	if srv == nil || srv.httpServer == nil {
		t.Fatal("Failed to initialize server")
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:18080/health")
	if err != nil {
		t.Fatalf("Failed to make request to server: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Server shutdown failed: %v", err)
	}

	startErr := <-errCh
	if startErr != http.ErrServerClosed {
		t.Errorf("Expected Start() to return ErrServerClosed, got %v", startErr)
	}
}

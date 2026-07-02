package httptelemetry_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kishan-thanki/logger/v2/httptelemetry"
	"github.com/kishan-thanki/logger/v2/slogctx"
)

func TestMiddleware(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	ctxHandler := slogctx.NewHandler(baseHandler)
	log := slog.New(ctxHandler)
	slog.SetDefault(log)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "Inside handler")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK"))
	})

	middlewareHandler := httptelemetry.Middleware(nextHandler)

	req := httptest.NewRequest("GET", "/test-path", nil)
	req.Header.Set("X-Trace-ID", "test-trace-123")
	rr := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != 2 {
		t.Fatalf("Expected 2 log lines, got %d", len(lines))
	}

	var log1 map[string]interface{}
	if err := json.Unmarshal(lines[0], &log1); err != nil {
		t.Fatalf("Failed to unmarshal first log line: %v", err)
	}
	if log1["msg"] != "Inside handler" {
		t.Errorf("Expected 'Inside handler', got %v", log1["msg"])
	}
	if log1["trace_id"] != "test-trace-123" {
		t.Errorf("Expected 'test-trace-123', got %v", log1["trace_id"])
	}

	var log2 map[string]interface{}
	if err := json.Unmarshal(lines[1], &log2); err != nil {
		t.Fatalf("Failed to unmarshal second log line: %v", err)
	}
	if log2["msg"] != "HTTP Request" {
		t.Errorf("Expected 'HTTP Request', got %v", log2["msg"])
	}
	if log2["trace_id"] != "test-trace-123" {
		t.Errorf("Expected 'test-trace-123', got %v", log2["trace_id"])
	}
	responseMap, ok := log2["response"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected response map, got nil")
	}
	if status, ok := responseMap["status"].(float64); !ok || int(status) != http.StatusCreated {
		t.Errorf("Expected status 201, got %v", responseMap["status"])
	}
}

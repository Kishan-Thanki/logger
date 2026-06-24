package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/kishan-thanki/logger"
)

func TestCoreHandler_PIIRedaction(t *testing.T) {
	var buf bytes.Buffer

	log := logger.New(
		logger.WithOutput(&buf),
		logger.WithRedaction("password", "secret_token", "credit_card"),
		logger.WithTraceID(false),
	)

	log.Info("User login attempt",
		slog.String("username", "admin"),
		slog.String("password", "mySuperSecretPassword123"),
		slog.String("credit_card", "4111-1111-1111-1111"),
	)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result["password"] != "***REDACTED***" {
		t.Errorf("Expected password to be redacted, got %v", result["password"])
	}
	if result["credit_card"] != "***REDACTED***" {
		t.Errorf("Expected credit_card to be redacted, got %v", result["credit_card"])
	}

	if result["username"] != "admin" {
		t.Errorf("Expected username to be 'admin', got %v", result["username"])
	}
}

func TestCoreHandler_TraceID(t *testing.T) {
	var buf bytes.Buffer

	log := logger.New(
		logger.WithOutput(&buf),
		logger.WithTraceID(true),
	)

	ctx := logger.InjectTraceID(context.Background(), "trace-xyz-123")

	log.InfoContext(ctx, "Processing payment")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result["trace_id"] != "trace-xyz-123" {
		t.Errorf("Expected trace_id to be 'trace-xyz-123', got %v", result["trace_id"])
	}
}

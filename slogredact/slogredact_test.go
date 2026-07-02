package slogredact_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/kishan-thanki/logger/v2/slogredact"
)

func TestHandler_Redact(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := slogredact.NewHandler(baseHandler, "password", "token")
	log := slog.New(handler)

	log.InfoContext(context.Background(), "Login attempt",
		slog.String("username", "admin"),
		slog.String("password", "supersecret123"),
		slog.String("token", "xyz789"),
	)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result["password"] != "***REDACTED***" {
		t.Errorf("Expected password to be redacted, got %v", result["password"])
	}
	if result["token"] != "***REDACTED***" {
		t.Errorf("Expected token to be redacted, got %v", result["token"])
	}
	if result["username"] != "admin" {
		t.Errorf("Expected username to be 'admin', got %v", result["username"])
	}
}

func TestReplaceAttr(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		ReplaceAttr: slogredact.ReplaceAttr("secret"),
	})
	log := slog.New(baseHandler)

	log.Info("Action", slog.String("secret", "my-secret"))

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result["secret"] != "***REDACTED***" {
		t.Errorf("Expected secret to be redacted, got %v", result["secret"])
	}
}

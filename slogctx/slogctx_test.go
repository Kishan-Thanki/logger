package slogctx_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/kishan-thanki/logger/slogctx"
)

func TestHandler_InjectTraceID(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := slogctx.NewHandler(baseHandler)
	log := slog.New(handler)

	ctx := slogctx.InjectTraceID(context.Background(), "REQ-999")

	log.InfoContext(ctx, "test message", slog.String("foo", "bar"))

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if result["trace_id"] != "REQ-999" {
		t.Errorf("Expected trace_id to be 'REQ-999', got %v", result["trace_id"])
	}

	if result["msg"] != "test message" {
		t.Errorf("Expected msg to be 'test message', got %v", result["msg"])
	}
	if result["foo"] != "bar" {
		t.Errorf("Expected foo to be 'bar', got %v", result["foo"])
	}
}

func TestHandler_NoTraceID(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := slogctx.NewHandler(baseHandler)
	log := slog.New(handler)

	log.InfoContext(context.Background(), "test message")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if _, exists := result["trace_id"]; exists {
		t.Errorf("Expected no trace_id in output, but found one: %v", result["trace_id"])
	}
}

func TestHandler_WithAttrsAndGroup(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := slogctx.NewHandler(baseHandler)
	log := slog.New(handler)

	ctx := slogctx.InjectTraceID(context.Background(), "REQ-777")

	childLog := log.With("env", "test").WithGroup("system")
	childLog.InfoContext(ctx, "system starting", "cpu", 100)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	system, ok := result["system"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'system' group to be a map, got %T", result["system"])
	}

	if system["trace_id"] != "REQ-777" {
		t.Errorf("Expected trace_id to be 'REQ-777' inside system group, got %v", system["trace_id"])
	}

	if result["env"] != "test" {
		t.Errorf("Expected env to be 'test', got %v", result["env"])
	}

	if system["cpu"] != float64(100) {
		t.Errorf("Expected system.cpu to be 100, got %v", system["cpu"])
	}
}

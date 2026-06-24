package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// coreHandler wraps a standard slog.Handler to provide automatic
// context extraction (Trace IDs). PII redaction is handled via ReplaceAttr.
type coreHandler struct {
	slog.Handler
	traceEnabled bool
}

// Handle intercepts the log record, injects context telemetry, and forwards it.
func (h *coreHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.traceEnabled {
		if traceID := ExtractTraceID(ctx); traceID != "" {
			r.AddAttrs(slog.String("trace_id", traceID))
		}
	}
	return h.Handler.Handle(ctx, r)
}

// WithAttrs ensures the wrapper propagates properly when child loggers are created.
func (h *coreHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &coreHandler{
		Handler:      h.Handler.WithAttrs(attrs),
		traceEnabled: h.traceEnabled,
	}
}

// WithGroup ensures the wrapper propagates properly when groups are created.
func (h *coreHandler) WithGroup(name string) slog.Handler {
	return &coreHandler{
		Handler:      h.Handler.WithGroup(name),
		traceEnabled: h.traceEnabled,
	}
}

// New creates a new highly-optimized slog.Logger configured for production usage.
func New(opts ...Option) *slog.Logger {
	cfg := &Config{
		Level:      new(slog.LevelVar), // Defaults to INFO automatically
		RedactKeys: []string{},
		TraceID:    false,
		Output:     os.Stdout,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	redactMap := make(map[string]bool, len(cfg.RedactKeys))
	for _, k := range cfg.RedactKeys {
		redactMap[strings.ToLower(k)] = true
	}

	handlerOpts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.Source,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Redact sensitive keys
			if redactMap[strings.ToLower(a.Key)] {
				a.Value = slog.StringValue("***REDACTED***")
			}
			return a
		},
	}

	baseHandler := slog.NewJSONHandler(cfg.Output, handlerOpts)

	eh := &coreHandler{
		Handler:      baseHandler,
		traceEnabled: cfg.TraceID,
	}

	return slog.New(eh)
}

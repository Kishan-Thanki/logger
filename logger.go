package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type coreHandler struct {
	slog.Handler
	traceEnabled bool
}

func (h *coreHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.traceEnabled {
		if traceID := ExtractTraceID(ctx); traceID != "" {
			r.AddAttrs(slog.String("trace_id", traceID))
		}
	}
	return h.Handler.Handle(ctx, r)
}

func (h *coreHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &coreHandler{
		Handler:      h.Handler.WithAttrs(attrs),
		traceEnabled: h.traceEnabled,
	}
}

func (h *coreHandler) WithGroup(name string) slog.Handler {
	return &coreHandler{
		Handler:      h.Handler.WithGroup(name),
		traceEnabled: h.traceEnabled,
	}
}

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

	handlerOpts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.Source,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Redact sensitive keys
			for _, k := range cfg.RedactKeys {
				if strings.EqualFold(a.Key, k) {
					a.Value = slog.StringValue("***REDACTED***")
					break
				}
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

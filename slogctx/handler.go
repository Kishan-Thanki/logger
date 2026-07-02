package slogctx

import (
	"context"
	"log/slog"
)

type contextKey string

const traceIDKey contextKey = "trace_id"

func InjectTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func ExtractTraceID(ctx context.Context) string {
	if val, ok := ctx.Value(traceIDKey).(string); ok {
		return val
	}
	return ""
}

type Handler struct {
	slog.Handler
}

func NewHandler(next slog.Handler) *Handler {
	return &Handler{Handler: next}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if traceID := ExtractTraceID(ctx); traceID != "" {
		r.AddAttrs(slog.String("trace_id", traceID))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{Handler: h.Handler.WithGroup(name)}
}

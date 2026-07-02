package slogredact

import (
	"context"
	"log/slog"
	"strings"
)

func ReplaceAttr(keys ...string) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		for _, k := range keys {
			if strings.EqualFold(a.Key, k) {
				a.Value = slog.StringValue("***REDACTED***")
				break
			}
		}
		return a
	}
}

type Handler struct {
	next slog.Handler
	keys []string
}

func NewHandler(next slog.Handler, keys ...string) *Handler {
	return &Handler{
		next: next,
		keys: keys,
	}
}

func (h *Handler) redactAttr(a slog.Attr) slog.Attr {
	for _, k := range h.keys {
		if strings.EqualFold(a.Key, k) {
			return slog.String(a.Key, "***REDACTED***")
		}
	}

	if a.Value.Kind() == slog.KindGroup {
		attrs := a.Value.Group()
		for i, attr := range attrs {
			attrs[i] = h.redactAttr(attr)
		}

		anyAttrs := make([]any, len(attrs))
		for i, v := range attrs {
			anyAttrs[i] = v
		}
		return slog.Group(a.Key, anyAttrs...)
	}

	return a
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	var attrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, h.redactAttr(a))
		return true
	})

	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	newRecord.AddAttrs(attrs...)

	return h.next.Handle(ctx, newRecord)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	redacted := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		redacted[i] = h.redactAttr(a)
	}
	return &Handler{
		next: h.next.WithAttrs(redacted),
		keys: h.keys,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		next: h.next.WithGroup(name),
		keys: h.keys,
	}
}

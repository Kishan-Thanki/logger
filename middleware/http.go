package middleware

import (
	"encoding/binary"
	"encoding/hex"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/kishan-thanki/logger"
)

func generateTraceID() string {
	var b [16]byte
	binary.LittleEndian.PutUint64(b[0:8], rand.Uint64())
	binary.LittleEndian.PutUint64(b[8:16], rand.Uint64())
	return hex.EncodeToString(b[:])
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// HTTP logs incoming HTTP requests as a structured slog.Group and automatically injects a Trace ID.
func HTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = generateTraceID()
		}

		ctx := logger.InjectTraceID(r.Context(), traceID)
		r = r.WithContext(ctx)

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		slog.InfoContext(ctx, "HTTP Request",
			slog.Group("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("ip", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			),
			slog.Group("response",
				slog.Int("status", rw.status),
				slog.String("latency", duration.String()),
				slog.Int64("latency_ms", duration.Milliseconds()),
			),
		)
	})
}

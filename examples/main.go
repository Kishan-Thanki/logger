package main

import (
	"log/slog"
	"net/http"

	"github.com/kishan-thanki/logger"
	"github.com/kishan-thanki/logger/middleware"
)

func main() {
	// 1. Initialize the core slog Handler
	handler := logger.New(
		logger.WithLevel("INFO"),
		logger.WithTraceID(true),
		logger.WithRedaction("password", "token", "credit_card"),
	)

	// 2. Set it as the Go global default logger
	slog.SetDefault(handler)

	// 3. Set up a simple HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		// Log the attempt (password will be automatically redacted if present)
		slog.InfoContext(r.Context(), "Processing login request",
			slog.String("username", "admin_user"),
			slog.String("password", "super_secret_password"),
		)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	})

	// 4. Wrap the mux with our telemetry middleware
	// This automatically injects Trace IDs into the context and logs all HTTP metrics!
	loggedMux := middleware.HTTP(mux)

	slog.Info("Starting server", slog.Int("port", 8080))
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		slog.Error("Server failed", slog.Any("error", err))
	}
}

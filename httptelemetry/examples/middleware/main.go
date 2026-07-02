package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/kishan-thanki/logger/v2/httptelemetry"
	"github.com/kishan-thanki/logger/v2/slogctx"
)

var base = slog.NewJSONHandler(os.Stdout, nil)
var log = slog.New(slogctx.NewHandler(base))

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.InfoContext(r.Context(),
		"home endpoint called",
		"method", r.Method,
		"path", r.URL.Path,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome!",
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	log.InfoContext(r.Context(), "health check")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.ErrorContext(r.Context(),
			"invalid login request",
			"error", err,
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.InfoContext(r.Context(),
		"login attempt",
		"email", req.Email,
	)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Login successful",
		"user": map[string]string{
			"email": req.Email,
		},
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /login", loginHandler)

	handler := httptelemetry.Middleware(mux)

	log.Info("starting server", "addr", ":8080")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Error("server failed", "error", err)
	}
}

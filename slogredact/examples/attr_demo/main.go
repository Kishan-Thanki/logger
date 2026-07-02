package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/kishan-thanki/logger/v2/slogredact"
)

// Third-party handler.
// We can't pass HandlerOptions here.
// var log = slog.New(myhandler.NewHandler(os.Stdout))

// === WIRING UP OUR HANDLER INSTEAD ===
var myhandler = slog.NewJSONHandler(os.Stdout, nil)
var safeHandler = slogredact.NewHandler(myhandler, "password")
var middlewareLog = slog.New(safeHandler)

type MiddlewareLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func middlewareHomeHandler(w http.ResponseWriter, r *http.Request) {
	middlewareLog.Info("home endpoint called",
		"method", r.Method,
		"path", r.URL.Path,
	)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome!",
	})
}

func middlewareHealthHandler(w http.ResponseWriter, r *http.Request) {
	middlewareLog.Info("health check")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func middlewareLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req MiddlewareLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middlewareLog.Error("invalid login request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	middlewareLog.Info("login attempt",
		"email", req.Email,
		"password", req.Password,
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

	mux.HandleFunc("GET /", middlewareHomeHandler)
	mux.HandleFunc("GET /health", middlewareHealthHandler)
	mux.HandleFunc("POST /login", middlewareLoginHandler)

	middlewareLog.Info("starting server", "addr", ":8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		middlewareLog.Error("server failed", "error", err)
	}
}

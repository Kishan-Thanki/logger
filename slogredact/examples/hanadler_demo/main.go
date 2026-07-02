package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/kishan-thanki/logger/slogredact"
)

var baseHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	ReplaceAttr: slogredact.ReplaceAttr("email", "password"),
})
var attrLog = slog.New(baseHandler)

type AttrLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func attrHomeHandler(w http.ResponseWriter, r *http.Request) {
	attrLog.Info("home endpoint called",
		"method", r.Method,
		"path", r.URL.Path,
	)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome!",
	})
}

func attrHealthHandler(w http.ResponseWriter, r *http.Request) {
	attrLog.Info("health check")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func attrLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req AttrLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		attrLog.Error("invalid login request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	attrLog.Info("login attempt",
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

	mux.HandleFunc("GET /", attrHomeHandler)
	mux.HandleFunc("GET /health", attrHealthHandler)
	mux.HandleFunc("POST /login", attrLoginHandler)

	attrLog.Info("starting server", "addr", ":8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		attrLog.Error("server failed", "error", err)
	}
}

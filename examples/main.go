package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/kishan-thanki/logger/v2/httptelemetry"
	"github.com/kishan-thanki/logger/v2/slogctx"
	"github.com/kishan-thanki/logger/v2/slogredact"
)

var base = slog.NewJSONHandler(os.Stdout, nil)
var safeHandler = slogredact.NewHandler(base, "password", "token")
var ctxHandler = slogctx.NewHandler(safeHandler)
var log = slog.New(ctxHandler)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.ErrorContext(r.Context(), "invalid login request", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.InfoContext(r.Context(), "login attempt",
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
	mux.HandleFunc("POST /login", loginHandler)

	handler := httptelemetry.Middleware(mux)

	log.Info("starting server", "addr", ":8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Error("server failed", "error", err)
	}
}

package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"tabler-api/auth"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
        "service":   "tabler-api",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// listUsersHandler returns a simple list of users as JSON.
func listUsersHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Authenticate request
    _, err := auth.UserFromRequest(r)
    if err != nil {
        slog.Error("failed to get user", slog.Any("error", err))
        w.WriteHeader(http.StatusUnauthorized)
        _ = json.NewEncoder(w).Encode(ErrorResponse{Error: "Authentication failed"})
        return
    }

    users := []auth.User{
        {ID: "1", Email: "alice@example.com", Name: "Alice"},
        {ID: "2", Email: "bob@example.com", Name: "Bob"},
    }

    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(users)
}

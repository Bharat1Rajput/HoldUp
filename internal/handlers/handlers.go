package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
	}

	writeJSON(w, http.StatusOK, response)
}

// ResourceHandler simulates a protected API resource.
func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message":   "resource accessed successfully",
		"timestamp": time.Now().UTC(),
	}

	writeJSON(w, http.StatusOK, response)
}

// writeJSON is a small helper to standardize JSON responses.
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

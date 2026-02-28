package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// HandlePosition returns an http.HandlerFunc that serves the current ISS position as JSON.
func HandlePosition(cache *TLECache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, l1, l2, err := cache.Get()
		if err != nil {
			slog.Error("failed to get TLE", "error", err)
			writeError(w, "TLE_UNAVAILABLE", "TLE unavailable", http.StatusServiceUnavailable)
			return
		}

		pos, err := GetCurrentPosition(name, l1, l2)
		if err != nil {
			slog.Error("failed to compute position", "error", err)
			writeError(w, "POSITION_ERROR", "position calculation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(pos); err != nil {
			slog.Error("failed to encode position", "error", err)
		}
	}
}

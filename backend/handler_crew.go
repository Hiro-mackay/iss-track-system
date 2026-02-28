package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// HandleCrew returns an http.HandlerFunc that serves the current ISS crew as JSON.
func HandleCrew(crew *CrewCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		members, err := crew.Get()
		if err != nil {
			slog.Error("failed to get crew", "error", err)
			writeError(w, "CREW_UNAVAILABLE", "crew data unavailable", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(members); err != nil {
			slog.Error("failed to encode crew", "error", err)
		}
	}
}

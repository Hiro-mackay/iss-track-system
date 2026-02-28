package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// HandleStatus returns an http.HandlerFunc that serves ISS status information as JSON.
func HandleStatus(tle *TLECache, crew *CrewCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, l1, l2, err := tle.Get()
		if err != nil {
			slog.Error("failed to get TLE", "error", err)
			writeError(w, "TLE_UNAVAILABLE", "TLE unavailable", http.StatusServiceUnavailable)
			return
		}

		crewCount := 0
		if members, err := crew.Get(); err == nil {
			crewCount = len(members)
		}

		status, err := ComputeISSStatus(name, l1, l2, crewCount)
		if err != nil {
			slog.Error("failed to compute ISS status", "error", err)
			writeError(w, "STATUS_ERROR", "status calculation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			slog.Error("failed to encode status", "error", err)
		}
	}
}

package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

// HandleOrbit returns an http.HandlerFunc that serves the ISS orbit track as JSON.
func HandleOrbit(cache *TLECache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minutes := 90
		if q := r.URL.Query().Get("minutes"); q != "" {
			if v, err := strconv.Atoi(q); err == nil && v > 0 {
				minutes = v
			}
		}
		if minutes > 180 {
			minutes = 180
		}

		name, l1, l2, err := cache.Get()
		if err != nil {
			slog.Error("failed to get TLE", "error", err)
			writeError(w, "TLE_UNAVAILABLE", "TLE unavailable", http.StatusServiceUnavailable)
			return
		}

		points, err := GetOrbitTrack(name, l1, l2, minutes)
		if err != nil {
			slog.Error("failed to compute orbit", "error", err)
			writeError(w, "ORBIT_ERROR", "orbit calculation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(points); err != nil {
			slog.Error("failed to encode orbit", "error", err)
		}
	}
}

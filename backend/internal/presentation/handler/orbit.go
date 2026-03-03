package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/presentation"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/service"
	trackingsvc "github.com/Hiro-mackay/iss-track-system/backend/internal/service/tracking"
)

// Orbit returns an http.HandlerFunc that serves the ISS orbit track as JSON.
func Orbit(svc *trackingsvc.QueryService) http.HandlerFunc {
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

		points, err := svc.OrbitTrack(minutes)
		if err != nil {
			var tleErr *service.TLEError
			if errors.As(err, &tleErr) {
				slog.Error("failed to get TLE", "error", err)
				presentation.WriteError(w, "TLE_UNAVAILABLE", "TLE unavailable", http.StatusServiceUnavailable)
				return
			}
			slog.Error("failed to compute orbit", "error", err)
			presentation.WriteError(w, "ORBIT_ERROR", "orbit calculation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(points); err != nil {
			slog.Error("failed to encode orbit", "error", err)
		}
	}
}

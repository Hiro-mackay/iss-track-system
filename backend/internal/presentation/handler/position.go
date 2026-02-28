package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/presentation"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/service/query"
)

// Position returns an http.HandlerFunc that serves the current ISS position as JSON.
func Position(svc *query.PositionQueryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pos, err := svc.CurrentPosition()
		if err != nil {
			var tleErr *query.TLEError
			if errors.As(err, &tleErr) {
				slog.Error("failed to get TLE", "error", err)
				presentation.WriteError(w, "TLE_UNAVAILABLE", "TLE unavailable", http.StatusServiceUnavailable)
				return
			}
			slog.Error("failed to compute position", "error", err)
			presentation.WriteError(w, "POSITION_ERROR", "position calculation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(pos); err != nil {
			slog.Error("failed to encode position", "error", err)
		}
	}
}

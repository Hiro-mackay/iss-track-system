package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/presentation"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/service/query"
)

// Status returns an http.HandlerFunc that serves ISS status information as JSON.
func Status(svc *query.StationQueryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := svc.Status()
		if err != nil {
			var tleErr *query.TLEError
			if errors.As(err, &tleErr) {
				slog.Error("failed to get TLE", "error", err)
				presentation.WriteError(w, "TLE_UNAVAILABLE", "TLE unavailable", http.StatusServiceUnavailable)
				return
			}
			slog.Error("failed to compute ISS status", "error", err)
			presentation.WriteError(w, "STATUS_ERROR", "status calculation failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			slog.Error("failed to encode status", "error", err)
		}
	}
}

package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/presentation"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/service/query"
)

// Passes returns an http.HandlerFunc that serves ISS pass predictions as JSON.
func Passes(svc *query.PassQueryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		latStr := r.URL.Query().Get("lat")
		lonStr := r.URL.Query().Get("lon")
		daysStr := r.URL.Query().Get("days")

		if latStr == "" || lonStr == "" {
			presentation.WriteError(w, "VALIDATION_ERROR", "lat and lon are required", http.StatusBadRequest)
			return
		}

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil || lat < -90 || lat > 90 {
			presentation.WriteError(w, "VALIDATION_ERROR", "invalid lat: must be between -90 and 90", http.StatusBadRequest)
			return
		}

		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil || lon < -180 || lon > 180 {
			presentation.WriteError(w, "VALIDATION_ERROR", "invalid lon: must be between -180 and 180", http.StatusBadRequest)
			return
		}

		days := 3
		if daysStr != "" {
			d, err := strconv.Atoi(daysStr)
			if err != nil || d < 1 || d > 7 {
				presentation.WriteError(w, "VALIDATION_ERROR", "invalid days: must be between 1 and 7", http.StatusBadRequest)
				return
			}
			days = d
		}

		passes, err := svc.Predict(lat, lon, days)
		if err != nil {
			var tleErr *query.TLEError
			if errors.As(err, &tleErr) {
				slog.Error("failed to get TLE", "error", err)
				presentation.WriteError(w, "TLE_UNAVAILABLE", "TLE unavailable", http.StatusServiceUnavailable)
				return
			}
			slog.Error("failed to compute pass predictions", "error", err)
			presentation.WriteError(w, "PASS_PREDICTION_ERROR", "pass prediction failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(passes); err != nil {
			slog.Error("failed to encode passes", "error", err)
		}
	}
}

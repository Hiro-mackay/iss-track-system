package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/presentation"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/service/query"
)

// Crew returns an http.HandlerFunc that serves the current ISS crew as JSON.
func Crew(svc *query.CrewQueryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		members, err := svc.List()
		if err != nil {
			slog.Error("failed to get crew", "error", err)
			presentation.WriteError(w, "CREW_UNAVAILABLE", "crew data unavailable", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(members); err != nil {
			slog.Error("failed to encode crew", "error", err)
		}
	}
}

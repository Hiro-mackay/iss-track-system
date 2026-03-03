package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/config"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/presentation/handler"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/presentation/middleware"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/repository"
	passsvc "github.com/Hiro-mackay/iss-track-system/backend/internal/service/pass"
	stationsvc "github.com/Hiro-mackay/iss-track-system/backend/internal/service/station"
	trackingsvc "github.com/Hiro-mackay/iss-track-system/backend/internal/service/tracking"
)

func main() {
	cfg := config.Load()

	// Repositories
	tleRepo := repository.NewTLERepository(cfg)
	if err := tleRepo.Refresh(); err != nil {
		slog.Error("initial TLE refresh failed", "error", err)
	}

	go func() {
		ticker := time.NewTicker(cfg.TLECacheTTL)
		defer ticker.Stop()
		for range ticker.C {
			if err := tleRepo.Refresh(); err != nil {
				slog.Error("periodic TLE refresh failed", "error", err)
			}
		}
	}()

	crewRepo := repository.NewCrewRepository(cfg)
	if err := crewRepo.Refresh(); err != nil {
		slog.Error("initial crew refresh failed", "error", err)
	}

	go func() {
		ticker := time.NewTicker(cfg.CrewCacheTTL)
		defer ticker.Stop()
		for range ticker.C {
			if err := crewRepo.Refresh(); err != nil {
				slog.Error("periodic crew refresh failed", "error", err)
			}
		}
	}()

	// Query services
	positionSvc := trackingsvc.NewQueryService(tleRepo)
	passSvc := passsvc.NewQueryService(tleRepo)
	stationSvc := stationsvc.NewQueryService(tleRepo, crewRepo)

	// Routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !tleRepo.HasData() {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"degraded","reason":"TLE cache empty"}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("GET /api/v1/position", handler.Position(positionSvc))
	mux.HandleFunc("GET /api/v1/orbit", handler.Orbit(positionSvc))
	mux.HandleFunc("GET /api/v1/passes", handler.Passes(passSvc))
	mux.HandleFunc("GET /api/v1/iss/crew", handler.Crew(stationSvc))
	mux.HandleFunc("GET /api/v1/iss/status", handler.Status(stationSvc))
	mux.Handle("GET /ws/position", handler.WebSocket(positionSvc, cfg.WSInterval, cfg.CORSOrigins))

	h := middleware.CORS(cfg.CORSOrigins, mux)

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := LoadConfig()

	cache := NewTLECache(cfg)
	if err := cache.Refresh(); err != nil {
		slog.Error("initial TLE refresh failed", "error", err)
	}

	// Periodic TLE refresh.
	go func() {
		ticker := time.NewTicker(cfg.TLECacheTTL)
		defer ticker.Stop()
		for range ticker.C {
			if err := cache.Refresh(); err != nil {
				slog.Error("periodic TLE refresh failed", "error", err)
			}
		}
	}()

	crewCache := NewCrewCache(cfg)
	if err := crewCache.Refresh(); err != nil {
		slog.Error("initial crew refresh failed", "error", err)
	}

	// Periodic crew refresh.
	go func() {
		ticker := time.NewTicker(cfg.CrewCacheTTL)
		defer ticker.Stop()
		for range ticker.C {
			if err := crewCache.Refresh(); err != nil {
				slog.Error("periodic crew refresh failed", "error", err)
			}
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if !cache.HasData() {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"degraded","reason":"TLE cache empty"}`))
			return
		}
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("GET /api/v1/position", HandlePosition(cache))
	mux.HandleFunc("GET /api/v1/orbit", HandleOrbit(cache))
	mux.HandleFunc("GET /api/v1/passes", HandlePasses(cache))
	mux.HandleFunc("GET /api/v1/iss/crew", HandleCrew(crewCache))
	mux.HandleFunc("GET /api/v1/iss/status", HandleStatus(cache, crewCache))
	mux.Handle("GET /ws/position", HandleWebSocket(cache, cfg.WSInterval, cfg.CORSOrigins))

	handler := CORSMiddleware(cfg.CORSOrigins, mux)

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           handler,
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

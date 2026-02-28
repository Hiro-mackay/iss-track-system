package repository

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/config"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/model"
)

var _ model.TLEProvider = (*TLERepository)(nil)

// tleMaxResponseBytes limits the size of TLE responses to 64 KiB.
const tleMaxResponseBytes = 1 << 16

// TLERepository provides thread-safe caching of TLE data with automatic refresh.
type TLERepository struct {
	mu        sync.RWMutex
	refreshMu sync.Mutex
	name      string
	line1     string
	line2     string
	fetchedAt time.Time
	cfg       config.Config
	client    *http.Client
}

// NewTLERepository creates a new TLERepository with the given configuration.
func NewTLERepository(cfg config.Config) *TLERepository {
	return &TLERepository{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Get returns cached TLE lines, refreshing if the cache has expired.
func (r *TLERepository) Get() (string, string, string, error) {
	r.mu.RLock()
	if r.name != "" && time.Since(r.fetchedAt) < r.cfg.TLECacheTTL {
		name, l1, l2 := r.name, r.line1, r.line2
		r.mu.RUnlock()
		return name, l1, l2, nil
	}
	r.mu.RUnlock()

	r.refreshMu.Lock()
	defer r.refreshMu.Unlock()

	// Double-check: another goroutine may have refreshed while we waited.
	r.mu.RLock()
	if r.name != "" && time.Since(r.fetchedAt) < r.cfg.TLECacheTTL {
		name, l1, l2 := r.name, r.line1, r.line2
		r.mu.RUnlock()
		return name, l1, l2, nil
	}
	r.mu.RUnlock()

	if err := r.Refresh(); err != nil {
		r.mu.RLock()
		if r.name != "" {
			name, l1, l2 := r.name, r.line1, r.line2
			r.mu.RUnlock()
			slog.Warn("TLE refresh failed, returning stale cache", "error", err)
			return name, l1, l2, nil
		}
		r.mu.RUnlock()
		return "", "", "", err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.name, r.line1, r.line2, nil
}

// Refresh fetches fresh TLE data from the configured URL and updates the cache.
func (r *TLERepository) Refresh() error {
	resp, err := r.client.Get(r.cfg.TLEURL)
	if err != nil {
		return fmt.Errorf("fetching TLE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TLE fetch returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, tleMaxResponseBytes))
	if err != nil {
		return fmt.Errorf("reading TLE response: %w", err)
	}

	var lines []string
	for _, line := range strings.Split(strings.TrimSpace(string(body)), "\n") {
		if s := strings.TrimSpace(line); s != "" {
			lines = append(lines, s)
		}
	}
	if len(lines) < 3 {
		return fmt.Errorf("invalid TLE data: expected 3 lines, got %d", len(lines))
	}

	r.mu.Lock()
	r.name = lines[0]
	r.line1 = lines[1]
	r.line2 = lines[2]
	r.fetchedAt = time.Now()
	r.mu.Unlock()

	slog.Info("TLE cache refreshed", "name", lines[0])
	return nil
}

// HasData reports whether the cache contains any TLE data.
func (r *TLERepository) HasData() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.name != ""
}

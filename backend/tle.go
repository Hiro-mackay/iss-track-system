package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

// tleMaxResponseBytes limits the size of TLE responses to 64 KiB.
const tleMaxResponseBytes = 1 << 16

// TLECache provides thread-safe caching of TLE data with automatic refresh.
type TLECache struct {
	mu        sync.RWMutex
	refreshMu sync.Mutex
	name      string
	line1     string
	line2     string
	fetchedAt time.Time
	cfg       Config
	client    *http.Client
}

// NewTLECache creates a new TLECache with the given configuration.
func NewTLECache(cfg Config) *TLECache {
	return &TLECache{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Get returns cached TLE lines, refreshing if the cache has expired.
func (c *TLECache) Get() (string, string, string, error) {
	c.mu.RLock()
	if c.name != "" && time.Since(c.fetchedAt) < c.cfg.TLECacheTTL {
		name, l1, l2 := c.name, c.line1, c.line2
		c.mu.RUnlock()
		return name, l1, l2, nil
	}
	c.mu.RUnlock()

	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	// Double-check: another goroutine may have refreshed while we waited.
	c.mu.RLock()
	if c.name != "" && time.Since(c.fetchedAt) < c.cfg.TLECacheTTL {
		name, l1, l2 := c.name, c.line1, c.line2
		c.mu.RUnlock()
		return name, l1, l2, nil
	}
	c.mu.RUnlock()

	if err := c.Refresh(); err != nil {
		c.mu.RLock()
		if c.name != "" {
			name, l1, l2 := c.name, c.line1, c.line2
			c.mu.RUnlock()
			slog.Warn("TLE refresh failed, returning stale cache", "error", err)
			return name, l1, l2, nil
		}
		c.mu.RUnlock()
		return "", "", "", err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.name, c.line1, c.line2, nil
}

// Refresh fetches fresh TLE data from the configured URL and updates the cache.
func (c *TLECache) Refresh() error {
	resp, err := c.client.Get(c.cfg.TLEURL)
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

	c.mu.Lock()
	c.name = lines[0]
	c.line1 = lines[1]
	c.line2 = lines[2]
	c.fetchedAt = time.Now()
	c.mu.Unlock()

	slog.Info("TLE cache refreshed", "name", lines[0])
	return nil
}

// HasData reports whether the cache contains any TLE data.
func (c *TLECache) HasData() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.name != ""
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// crewMaxResponseBytes limits the size of crew responses to 64 KiB.
const crewMaxResponseBytes = 1 << 16

type openNotifyPerson struct {
	Name  string `json:"name"`
	Craft string `json:"craft"`
}

type openNotifyResponse struct {
	People  []openNotifyPerson `json:"people"`
	Number  int                `json:"number"`
	Message string             `json:"message"`
}

// CrewCache provides thread-safe caching of ISS crew data with automatic refresh.
type CrewCache struct {
	mu        sync.RWMutex
	refreshMu sync.Mutex
	members   []CrewMember
	fetchedAt time.Time
	cfg       Config
	client    *http.Client
}

// NewCrewCache creates a new CrewCache with the given configuration.
func NewCrewCache(cfg Config) *CrewCache {
	return &CrewCache{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Get returns cached crew members, refreshing if the cache has expired.
func (c *CrewCache) Get() ([]CrewMember, error) {
	c.mu.RLock()
	if c.members != nil && time.Since(c.fetchedAt) < c.cfg.CrewCacheTTL {
		result := make([]CrewMember, len(c.members))
		copy(result, c.members)
		c.mu.RUnlock()
		return result, nil
	}
	c.mu.RUnlock()

	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	// Double-check: another goroutine may have refreshed while we waited.
	c.mu.RLock()
	if c.members != nil && time.Since(c.fetchedAt) < c.cfg.CrewCacheTTL {
		result := make([]CrewMember, len(c.members))
		copy(result, c.members)
		c.mu.RUnlock()
		return result, nil
	}
	c.mu.RUnlock()

	if err := c.Refresh(); err != nil {
		c.mu.RLock()
		if c.members != nil {
			result := make([]CrewMember, len(c.members))
			copy(result, c.members)
			c.mu.RUnlock()
			slog.Warn("crew refresh failed, returning stale cache", "error", err)
			return result, nil
		}
		c.mu.RUnlock()
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]CrewMember, len(c.members))
	copy(result, c.members)
	return result, nil
}

// Refresh fetches fresh crew data from the configured URL and updates the cache.
func (c *CrewCache) Refresh() error {
	resp, err := c.client.Get(c.cfg.CrewURL)
	if err != nil {
		return fmt.Errorf("fetching crew: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("crew fetch returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, crewMaxResponseBytes))
	if err != nil {
		return fmt.Errorf("reading crew response: %w", err)
	}

	var data openNotifyResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("parsing crew JSON: %w", err)
	}

	var members []CrewMember
	for _, p := range data.People {
		if p.Craft == "ISS" {
			members = append(members, CrewMember(p))
		}
	}

	c.mu.Lock()
	c.members = members
	c.fetchedAt = time.Now()
	c.mu.Unlock()

	slog.Info("crew cache refreshed", "count", len(members))
	return nil
}

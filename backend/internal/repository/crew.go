package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/config"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/station"
)

var _ station.CrewProvider = (*CrewRepository)(nil)

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

// CrewRepository provides thread-safe caching of ISS crew data with automatic refresh.
type CrewRepository struct {
	mu        sync.RWMutex
	refreshMu sync.Mutex
	members   []station.CrewMember
	fetchedAt time.Time
	cfg       config.Config
	client    *http.Client
}

// NewCrewRepository creates a new CrewRepository with the given configuration.
func NewCrewRepository(cfg config.Config) *CrewRepository {
	return &CrewRepository{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Get returns cached crew members, refreshing if the cache has expired.
func (r *CrewRepository) Get() ([]station.CrewMember, error) {
	r.mu.RLock()
	if r.members != nil && time.Since(r.fetchedAt) < r.cfg.CrewCacheTTL {
		result := make([]station.CrewMember, len(r.members))
		copy(result, r.members)
		r.mu.RUnlock()
		return result, nil
	}
	r.mu.RUnlock()

	r.refreshMu.Lock()
	defer r.refreshMu.Unlock()

	// Double-check: another goroutine may have refreshed while we waited.
	r.mu.RLock()
	if r.members != nil && time.Since(r.fetchedAt) < r.cfg.CrewCacheTTL {
		result := make([]station.CrewMember, len(r.members))
		copy(result, r.members)
		r.mu.RUnlock()
		return result, nil
	}
	r.mu.RUnlock()

	if err := r.Refresh(); err != nil {
		r.mu.RLock()
		if r.members != nil {
			result := make([]station.CrewMember, len(r.members))
			copy(result, r.members)
			r.mu.RUnlock()
			slog.Warn("crew refresh failed, returning stale cache", "error", err)
			return result, nil
		}
		r.mu.RUnlock()
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]station.CrewMember, len(r.members))
	copy(result, r.members)
	return result, nil
}

// Refresh fetches fresh crew data from the configured URL and updates the cache.
func (r *CrewRepository) Refresh() error {
	resp, err := r.client.Get(r.cfg.CrewURL)
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

	var members []station.CrewMember
	for _, p := range data.People {
		if p.Craft == "ISS" {
			members = append(members, station.CrewMember{Name: p.Name, Craft: p.Craft})
		}
	}

	r.mu.Lock()
	r.members = members
	r.fetchedAt = time.Now()
	r.mu.Unlock()

	slog.Info("crew cache refreshed", "count", len(members))
	return nil
}

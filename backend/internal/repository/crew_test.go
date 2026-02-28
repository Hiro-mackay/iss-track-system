package repository

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/config"
)

func TestCrewRepository_FiltersISSOnly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"people": [
				{"name": "Oleg Kononenko", "craft": "ISS"},
				{"name": "Nikolai Chub", "craft": "ISS"},
				{"name": "Tang Hongbo", "craft": "Tiangong"}
			],
			"number": 3,
			"message": "success"
		}`))
	}))
	defer srv.Close()

	cfg := config.Config{CrewURL: srv.URL, CrewCacheTTL: 1 * time.Hour}
	repo := NewCrewRepository(cfg)

	members, err := repo.Get()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(members) != 2 {
		t.Fatalf("should filter to ISS crew only: got %d, want 2", len(members))
	}

	for _, m := range members {
		if m.Craft != "ISS" {
			t.Errorf("should only contain ISS crew: got craft %q", m.Craft)
		}
	}
}

func TestCrewRepository_ParsesNames(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"people": [{"name": "Alice", "craft": "ISS"}],
			"number": 1,
			"message": "success"
		}`))
	}))
	defer srv.Close()

	cfg := config.Config{CrewURL: srv.URL, CrewCacheTTL: 1 * time.Hour}
	repo := NewCrewRepository(cfg)

	members, err := repo.Get()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if members[0].Name != "Alice" {
		t.Errorf("should parse crew name: got %q, want %q", members[0].Name, "Alice")
	}
}

func TestCrewRepository_ReturnsStaleFallback(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls > 1 {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"people": [{"name": "Bob", "craft": "ISS"}],
			"number": 1,
			"message": "success"
		}`))
	}))
	defer srv.Close()

	cfg := config.Config{CrewURL: srv.URL, CrewCacheTTL: 0}
	repo := NewCrewRepository(cfg)

	// First call succeeds
	members, err := repo.Get()
	if err != nil {
		t.Fatalf("first call should succeed: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("should have 1 member: got %d", len(members))
	}

	// Second call fails but returns stale data
	members, err = repo.Get()
	if err != nil {
		t.Fatalf("should return stale data on error: %v", err)
	}
	if len(members) != 1 || members[0].Name != "Bob" {
		t.Errorf("should return stale crew member Bob: got %v", members)
	}
}

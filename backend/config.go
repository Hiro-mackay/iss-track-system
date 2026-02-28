package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	Port         string
	TLEURL       string
	TLECacheTTL  time.Duration
	WSInterval   time.Duration
	CORSOrigins  []string
	CrewURL      string
	CrewCacheTTL time.Duration
}

// LoadConfig reads ISS_* environment variables and returns a Config with defaults.
func LoadConfig() Config {
	cfg := Config{
		Port:         ":8000",
		TLEURL:       "https://celestrak.org/NORAD/elements/gp.php?CATNR=25544&FORMAT=3LE",
		TLECacheTTL:  7200 * time.Second,
		WSInterval:   1 * time.Second,
		CORSOrigins:  []string{"http://localhost:3000"},
		CrewURL:      "http://api.open-notify.org/astros.json",
		CrewCacheTTL: 3600 * time.Second,
	}

	if v := os.Getenv("ISS_PORT"); v != "" {
		cfg.Port = v
	}

	if v := os.Getenv("ISS_TLE_URL"); v != "" {
		cfg.TLEURL = v
	}

	if v := os.Getenv("ISS_TLE_CACHE_TTL_SECONDS"); v != "" {
		if sec, err := strconv.Atoi(v); err == nil {
			cfg.TLECacheTTL = time.Duration(sec) * time.Second
		}
	}

	if v := os.Getenv("ISS_WS_INTERVAL_SECONDS"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.WSInterval = time.Duration(f * float64(time.Second))
		}
	}

	if v := os.Getenv("ISS_CORS_ORIGINS"); v != "" {
		var origins []string
		if err := json.Unmarshal([]byte(v), &origins); err == nil {
			cfg.CORSOrigins = origins
		}
	}

	if v := os.Getenv("ISS_CREW_URL"); v != "" {
		cfg.CrewURL = v
	}

	if v := os.Getenv("ISS_CREW_CACHE_TTL_SECONDS"); v != "" {
		if sec, err := strconv.Atoi(v); err == nil {
			cfg.CrewCacheTTL = time.Duration(sec) * time.Second
		}
	}

	return cfg
}

package main

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/akhenakh/sgp4"
)

// GetCurrentPosition calculates the current ISS position from TLE data.
func GetCurrentPosition(name, line1, line2 string) (ISSPosition, error) {
	tleStr := name + "\n" + line1 + "\n" + line2
	tle, err := sgp4.ParseTLE(tleStr)
	if err != nil {
		return ISSPosition{}, fmt.Errorf("parsing TLE: %w", err)
	}

	now := time.Now().UTC()
	eci, err := tle.FindPositionAtTime(now)
	if err != nil {
		return ISSPosition{}, fmt.Errorf("propagating position: %w", err)
	}

	lat, lon, alt := eci.ToGeodetic()

	vx := eci.Velocity.X
	vy := eci.Velocity.Y
	vz := eci.Velocity.Z
	speedKmh := math.Sqrt(vx*vx+vy*vy+vz*vz) * 3600

	return ISSPosition{
		Lat:         lat,
		Lon:         lon,
		AltitudeKm:  alt,
		VelocityKmh: speedKmh,
		Timestamp:   now.Format(time.RFC3339),
	}, nil
}

// GetOrbitTrack computes orbit points at 1-minute intervals centered on now.
func GetOrbitTrack(name, line1, line2 string, minutes int) ([]OrbitPoint, error) {
	tleStr := name + "\n" + line1 + "\n" + line2
	tle, err := sgp4.ParseTLE(tleStr)
	if err != nil {
		return nil, fmt.Errorf("parsing TLE: %w", err)
	}

	now := time.Now().UTC()
	start := now.Add(-time.Duration(minutes/2) * time.Minute)
	points := make([]OrbitPoint, 0, minutes)

	for i := 0; i < minutes; i++ {
		t := start.Add(time.Duration(i) * time.Minute)
		eci, err := tle.FindPositionAtTime(t)
		if err != nil {
			slog.Warn("orbit propagation failed", "offset_min", i, "error", err)
			continue
		}

		lat, lon, alt := eci.ToGeodetic()
		points = append(points, OrbitPoint{
			Lat:        lat,
			Lon:        lon,
			AltitudeKm: alt,
		})
	}

	if len(points) == 0 {
		return nil, fmt.Errorf("all orbit propagation points failed for %d minutes", minutes)
	}

	return points, nil
}

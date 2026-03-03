package tracking

import (
	"fmt"
	"log/slog"
	"math"
	"time"
)

// PropagatePosition calculates the current ISS position from the TLE data.
func (t TLE) PropagatePosition() (ISSPosition, error) {
	if !t.IsValid() {
		return ISSPosition{}, fmt.Errorf("invalid TLE: zero value")
	}

	now := time.Now().UTC()
	eci, err := t.parsed.FindPositionAtTime(now)
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

// PropagateOrbitTrack computes orbit points at 1-minute intervals centered on now.
func (t TLE) PropagateOrbitTrack(minutes int) ([]OrbitPoint, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("invalid TLE: zero value")
	}

	now := time.Now().UTC()
	start := now.Add(-time.Duration(minutes/2) * time.Minute)
	points := make([]OrbitPoint, 0, minutes)

	for i := 0; i < minutes; i++ {
		tm := start.Add(time.Duration(i) * time.Minute)
		eci, err := t.parsed.FindPositionAtTime(tm)
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

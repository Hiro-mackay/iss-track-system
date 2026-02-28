package orbit

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/model"
	"github.com/akhenakh/sgp4"
)

// PropagatePosition calculates the current ISS position from TLE data.
func PropagatePosition(name, line1, line2 string) (model.ISSPosition, error) {
	tleStr := name + "\n" + line1 + "\n" + line2
	tle, err := sgp4.ParseTLE(tleStr)
	if err != nil {
		return model.ISSPosition{}, fmt.Errorf("parsing TLE: %w", err)
	}

	now := time.Now().UTC()
	eci, err := tle.FindPositionAtTime(now)
	if err != nil {
		return model.ISSPosition{}, fmt.Errorf("propagating position: %w", err)
	}

	lat, lon, alt := eci.ToGeodetic()

	vx := eci.Velocity.X
	vy := eci.Velocity.Y
	vz := eci.Velocity.Z
	speedKmh := math.Sqrt(vx*vx+vy*vy+vz*vz) * 3600

	return model.ISSPosition{
		Lat:         lat,
		Lon:         lon,
		AltitudeKm:  alt,
		VelocityKmh: speedKmh,
		Timestamp:   now.Format(time.RFC3339),
	}, nil
}

// PropagateOrbitTrack computes orbit points at 1-minute intervals centered on now.
func PropagateOrbitTrack(name, line1, line2 string, minutes int) ([]model.OrbitPoint, error) {
	tleStr := name + "\n" + line1 + "\n" + line2
	tle, err := sgp4.ParseTLE(tleStr)
	if err != nil {
		return nil, fmt.Errorf("parsing TLE: %w", err)
	}

	now := time.Now().UTC()
	start := now.Add(-time.Duration(minutes/2) * time.Minute)
	points := make([]model.OrbitPoint, 0, minutes)

	for i := 0; i < minutes; i++ {
		t := start.Add(time.Duration(i) * time.Minute)
		eci, err := tle.FindPositionAtTime(t)
		if err != nil {
			slog.Warn("orbit propagation failed", "offset_min", i, "error", err)
			continue
		}

		lat, lon, alt := eci.ToGeodetic()
		points = append(points, model.OrbitPoint{
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

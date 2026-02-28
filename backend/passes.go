package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/akhenakh/sgp4"
)

// Brightness levels for pass predictions.
const (
	BrightnessBright   = "bright"
	BrightnessModerate = "moderate"
	BrightnessDim      = "dim"
)

// GetPassPredictions computes ISS pass predictions for an observer location.
func GetPassPredictions(name, line1, line2 string, lat, lon float64, days int) ([]PassPrediction, error) {
	tleStr := name + "\n" + line1 + "\n" + line2
	tle, err := sgp4.ParseTLE(tleStr)
	if err != nil {
		return nil, fmt.Errorf("parsing TLE: %w", err)
	}

	now := time.Now().UTC()
	stop := now.Add(time.Duration(days) * 24 * time.Hour)

	rawPasses, err := tle.GeneratePasses(lat, lon, 0, now, stop, 30)
	if err != nil {
		return nil, fmt.Errorf("generating passes: %w", err)
	}

	results := make([]PassPrediction, 0, len(rawPasses))
	for _, p := range rawPasses {
		if p.MaxElevation < 10 {
			continue
		}

		eci, err := tle.FindPositionAtTime(p.MaxElevationTime)
		if err != nil {
			slog.Warn("pass propagation failed", "time", p.MaxElevationTime, "error", err)
			continue
		}
		_, _, satAlt := eci.ToGeodetic()

		visible := IsPassVisible(tle, lat, lon, satAlt, p.MaxElevationTime)

		results = append(results, PassPrediction{
			Rise: PassEvent{
				Time:      p.AOS.Format(time.RFC3339),
				Azimuth:   p.AOSAzimuth,
				Elevation: 0,
			},
			Culmination: PassEvent{
				Time:      p.MaxElevationTime.Format(time.RFC3339),
				Azimuth:   p.MaxElevationAz,
				Elevation: p.MaxElevation,
			},
			Set: PassEvent{
				Time:      p.LOS.Format(time.RFC3339),
				Azimuth:   p.LOSAzimuth,
				Elevation: 0,
			},
			MaxElevation: p.MaxElevation,
			IsVisible:    visible,
			Brightness:   brightness(p.MaxElevation),
		})
	}

	return results, nil
}

func brightness(maxElev float64) string {
	switch {
	case maxElev >= 60:
		return BrightnessBright
	case maxElev >= 30:
		return BrightnessModerate
	default:
		return BrightnessDim
	}
}

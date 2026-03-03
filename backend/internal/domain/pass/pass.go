package pass

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/tracking"
)

// PredictPasses computes ISS pass predictions for an observer location.
func PredictPasses(tle tracking.TLE, lat, lon float64, days int) ([]PassPrediction, error) {
	if !tle.IsValid() {
		return nil, fmt.Errorf("invalid TLE: zero value")
	}

	parsed := tle.Parsed()
	now := time.Now().UTC()
	stop := now.Add(time.Duration(days) * 24 * time.Hour)

	rawPasses, err := parsed.GeneratePasses(lat, lon, 0, now, stop, 30)
	if err != nil {
		return nil, fmt.Errorf("generating passes: %w", err)
	}

	results := make([]PassPrediction, 0, len(rawPasses))
	for _, p := range rawPasses {
		if p.MaxElevation < 10 {
			continue
		}

		eci, err := parsed.FindPositionAtTime(p.MaxElevationTime)
		if err != nil {
			slog.Warn("pass propagation failed", "time", p.MaxElevationTime, "error", err)
			continue
		}
		_, _, satAlt := eci.ToGeodetic()

		visible := IsPassVisible(parsed, lat, lon, satAlt, p.MaxElevationTime)

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
			Brightness:   Brightness(p.MaxElevation),
		})
	}

	return results, nil
}

// Brightness returns the brightness classification for a given max elevation.
func Brightness(maxElev float64) string {
	switch {
	case maxElev >= 60:
		return BrightnessBright
	case maxElev >= 30:
		return BrightnessModerate
	default:
		return BrightnessDim
	}
}

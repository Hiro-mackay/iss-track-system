package orbit

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/model"
	"github.com/akhenakh/sgp4"
)

// PredictPasses computes ISS pass predictions for an observer location.
func PredictPasses(name, line1, line2 string, lat, lon float64, days int) ([]model.PassPrediction, error) {
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

	results := make([]model.PassPrediction, 0, len(rawPasses))
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

		results = append(results, model.PassPrediction{
			Rise: model.PassEvent{
				Time:      p.AOS.Format(time.RFC3339),
				Azimuth:   p.AOSAzimuth,
				Elevation: 0,
			},
			Culmination: model.PassEvent{
				Time:      p.MaxElevationTime.Format(time.RFC3339),
				Azimuth:   p.MaxElevationAz,
				Elevation: p.MaxElevation,
			},
			Set: model.PassEvent{
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
		return model.BrightnessBright
	case maxElev >= 30:
		return model.BrightnessModerate
	default:
		return model.BrightnessDim
	}
}

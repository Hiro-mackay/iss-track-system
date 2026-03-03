package station

import (
	"fmt"
	"math"
	"time"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/tracking"
)

// OrbitalParams holds the Keplerian orbital parameters for the ISS.
type OrbitalParams struct {
	InclinationDeg float64 `json:"inclination_deg"`
	PeriodMin      float64 `json:"period_min"`
	Eccentricity   float64 `json:"eccentricity"`
	ApogeeKm       float64 `json:"apogee_km"`
	PerigeeKm      float64 `json:"perigee_km"`
}

// ISSStatus holds static and computed information about the ISS.
type ISSStatus struct {
	LaunchYear    int           `json:"launch_year"`
	OrbitCount    int           `json:"orbit_count"`
	CrewCount     int           `json:"crew_count"`
	OrbitalParams OrbitalParams `json:"orbital_params"`
	TLEEpoch      string        `json:"tle_epoch"`
}

// sgp4EarthRadiusKm is the Earth equatorial radius used by SGP4 (WGS72).
const sgp4EarthRadiusKm = 6378.135

// ComputeStatus calculates ISS orbital status from TLE data and crew count.
func ComputeStatus(tle tracking.TLE, crewCount int) (ISSStatus, error) {
	if !tle.IsValid() {
		return ISSStatus{}, fmt.Errorf("invalid TLE: zero value")
	}

	parsed := tle.Parsed()
	periodMin := 1440.0 / parsed.MeanMotion
	smaER := parsed.RecoveredSemiMajorAxis()
	smaKm := smaER * sgp4EarthRadiusKm
	apogeeKm := smaKm*(1+parsed.Eccentricity) - sgp4EarthRadiusKm
	perigeeKm := smaKm*(1-parsed.Eccentricity) - sgp4EarthRadiusKm

	elapsed := time.Since(parsed.EpochTime()).Minutes()
	orbitCount := parsed.RevolutionNumber + int(elapsed*parsed.MeanMotion/1440.0)

	return ISSStatus{
		LaunchYear: 1998,
		OrbitCount: orbitCount,
		CrewCount:  crewCount,
		OrbitalParams: OrbitalParams{
			InclinationDeg: math.Round(parsed.Inclination*1000) / 1000,
			PeriodMin:      math.Round(periodMin*1000) / 1000,
			Eccentricity:   parsed.Eccentricity,
			ApogeeKm:       math.Round(apogeeKm*100) / 100,
			PerigeeKm:      math.Round(perigeeKm*100) / 100,
		},
		TLEEpoch: parsed.EpochTime().Format(time.RFC3339),
	}, nil
}

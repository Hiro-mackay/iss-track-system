package main

import (
	"fmt"
	"math"
	"time"

	"github.com/akhenakh/sgp4"
)

// sgp4EarthRadiusKm is the Earth equatorial radius used by SGP4 (WGS72).
const sgp4EarthRadiusKm = 6378.135

// ComputeISSStatus calculates ISS orbital status from TLE data and crew count.
func ComputeISSStatus(name, line1, line2 string, crewCount int) (ISSStatus, error) {
	tleStr := name + "\n" + line1 + "\n" + line2
	tle, err := sgp4.ParseTLE(tleStr)
	if err != nil {
		return ISSStatus{}, fmt.Errorf("parsing TLE: %w", err)
	}

	periodMin := 1440.0 / tle.MeanMotion
	smaER := tle.RecoveredSemiMajorAxis()
	smaKm := smaER * sgp4EarthRadiusKm
	apogeeKm := smaKm*(1+tle.Eccentricity) - sgp4EarthRadiusKm
	perigeeKm := smaKm*(1-tle.Eccentricity) - sgp4EarthRadiusKm

	elapsed := time.Since(tle.EpochTime()).Minutes()
	orbitCount := tle.RevolutionNumber + int(elapsed*tle.MeanMotion/1440.0)

	return ISSStatus{
		LaunchYear: 1998,
		OrbitCount: orbitCount,
		CrewCount:  crewCount,
		OrbitalParams: OrbitalParams{
			InclinationDeg: math.Round(tle.Inclination*1000) / 1000,
			PeriodMin:      math.Round(periodMin*1000) / 1000,
			Eccentricity:   tle.Eccentricity,
			ApogeeKm:       math.Round(apogeeKm*100) / 100,
			PerigeeKm:      math.Round(perigeeKm*100) / 100,
		},
		TLEEpoch: tle.EpochTime().Format(time.RFC3339),
	}, nil
}

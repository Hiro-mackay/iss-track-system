package model

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

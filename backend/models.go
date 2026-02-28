package main

// ISSPosition represents the current position and velocity of the ISS.
type ISSPosition struct {
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	AltitudeKm  float64 `json:"altitude_km"`
	VelocityKmh float64 `json:"velocity_kmh"`
	Timestamp   string  `json:"timestamp"`
}

// OrbitPoint represents a single point on the ISS orbit path.
type OrbitPoint struct {
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	AltitudeKm float64 `json:"altitude_km"`
}

// PassEvent represents a single event during a satellite pass (rise, culmination, or set).
type PassEvent struct {
	Time      string  `json:"time"`
	Azimuth   float64 `json:"azimuth"`
	Elevation float64 `json:"elevation"`
}

// PassPrediction represents a complete satellite pass over an observer location.
type PassPrediction struct {
	Rise         PassEvent `json:"rise"`
	Culmination  PassEvent `json:"culmination"`
	Set          PassEvent `json:"set"`
	MaxElevation float64   `json:"max_elevation"`
	IsVisible    bool      `json:"is_visible"`
	Brightness   string    `json:"brightness"`
}

// CrewMember represents a person currently aboard a spacecraft.
type CrewMember struct {
	Name  string `json:"name"`
	Craft string `json:"craft"`
}

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

package model

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

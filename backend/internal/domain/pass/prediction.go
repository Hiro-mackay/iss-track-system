package pass

// Brightness levels for pass predictions.
const (
	BrightnessBright   = "bright"
	BrightnessModerate = "moderate"
	BrightnessDim      = "dim"
)

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

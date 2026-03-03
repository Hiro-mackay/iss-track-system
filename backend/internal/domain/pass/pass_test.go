package pass

import (
	"testing"
	"time"
)

func TestBrightnessBright(t *testing.T) {
	if got := Brightness(60); got != "bright" {
		t.Errorf("Brightness(60) = %q, want %q", got, "bright")
	}
	if got := Brightness(90); got != "bright" {
		t.Errorf("Brightness(90) = %q, want %q", got, "bright")
	}
}

func TestBrightnessModerate(t *testing.T) {
	if got := Brightness(30); got != "moderate" {
		t.Errorf("Brightness(30) = %q, want %q", got, "moderate")
	}
	if got := Brightness(59); got != "moderate" {
		t.Errorf("Brightness(59) = %q, want %q", got, "moderate")
	}
}

func TestBrightnessDim(t *testing.T) {
	if got := Brightness(10); got != "dim" {
		t.Errorf("Brightness(10) = %q, want %q", got, "dim")
	}
	if got := Brightness(29); got != "dim" {
		t.Errorf("Brightness(29) = %q, want %q", got, "dim")
	}
}

func TestSunAltitudeRange(t *testing.T) {
	times := []time.Time{
		time.Date(2024, 6, 21, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 20, 6, 0, 0, 0, time.UTC),
	}
	locations := [][2]float64{
		{35.6762, 139.6503},  // Tokyo
		{51.5074, -0.1278},   // London
		{-33.8688, 151.2093}, // Sydney
		{0.0, 0.0},           // Equator/prime meridian
	}

	for _, loc := range locations {
		for _, tm := range times {
			alt := SunAltitude(loc[0], loc[1], tm)
			if alt < -90 || alt > 90 {
				t.Errorf("SunAltitude(%f, %f, %v) = %f, want in [-90, 90]",
					loc[0], loc[1], tm, alt)
			}
		}
	}
}

func TestSunAltitudeNoonHigher(t *testing.T) {
	noon := time.Date(2024, 3, 20, 12, 0, 0, 0, time.UTC)
	midnight := time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)

	noonAlt := SunAltitude(0, 0, noon)
	midnightAlt := SunAltitude(0, 0, midnight)

	if noonAlt <= midnightAlt {
		t.Errorf("noon alt (%f) should be > midnight alt (%f) at equator",
			noonAlt, midnightAlt)
	}
}

func TestPassPredictions(t *testing.T) {
	preds, err := PredictPasses(testTLE(), 35.6762, 139.6503, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, p := range preds {
		if p.MaxElevation < 10 {
			t.Errorf("pass[%d]: MaxElevation %f < 10", i, p.MaxElevation)
		}

		rise, err := time.Parse(time.RFC3339, p.Rise.Time)
		if err != nil {
			t.Errorf("pass[%d]: invalid rise time: %v", i, err)
			continue
		}
		culm, err := time.Parse(time.RFC3339, p.Culmination.Time)
		if err != nil {
			t.Errorf("pass[%d]: invalid culmination time: %v", i, err)
			continue
		}
		set, err := time.Parse(time.RFC3339, p.Set.Time)
		if err != nil {
			t.Errorf("pass[%d]: invalid set time: %v", i, err)
			continue
		}

		if !rise.Before(culm) {
			t.Errorf("pass[%d]: rise (%v) not before culmination (%v)", i, rise, culm)
		}
		if !culm.Before(set) {
			t.Errorf("pass[%d]: culmination (%v) not before set (%v)", i, culm, set)
		}

		switch p.Brightness {
		case "bright", "moderate", "dim":
		default:
			t.Errorf("pass[%d]: unexpected brightness %q", i, p.Brightness)
		}
	}
}

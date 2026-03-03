package tracking

import "testing"

func TestNewTLE_Valid(t *testing.T) {
	tle, err := NewTLE(testTLEName, testTLELine1, testTLELine2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tle.IsValid() {
		t.Error("TLE should be valid")
	}
	if tle.Name() != testTLEName {
		t.Errorf("Name() = %q, want %q", tle.Name(), testTLEName)
	}
}

func TestNewTLE_Invalid(t *testing.T) {
	_, err := NewTLE("bad", "bad", "bad")
	if err == nil {
		t.Error("should return error for invalid TLE")
	}
}

func TestTLEZeroValue(t *testing.T) {
	var tle TLE
	if tle.IsValid() {
		t.Error("zero-value TLE should not be valid")
	}
}

func TestPositionLatRange(t *testing.T) {
	pos, err := testTLE().PropagatePosition()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pos.Lat < -90 || pos.Lat > 90 {
		t.Errorf("lat %f out of range [-90, 90]", pos.Lat)
	}
}

func TestPositionLonRange(t *testing.T) {
	pos, err := testTLE().PropagatePosition()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pos.Lon < -180 || pos.Lon > 180 {
		t.Errorf("lon %f out of range [-180, 180]", pos.Lon)
	}
}

func TestPositionAltitudeRange(t *testing.T) {
	pos, err := testTLE().PropagatePosition()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pos.AltitudeKm < 200 || pos.AltitudeKm > 500 {
		t.Errorf("altitude %f out of range [200, 500]", pos.AltitudeKm)
	}
}

func TestPositionVelocityRange(t *testing.T) {
	pos, err := testTLE().PropagatePosition()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pos.VelocityKmh < 25000 || pos.VelocityKmh > 30000 {
		t.Errorf("velocity %f out of range [25000, 30000]", pos.VelocityKmh)
	}
}

func TestPositionTimestamp(t *testing.T) {
	pos, err := testTLE().PropagatePosition()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pos.Timestamp == "" {
		t.Error("timestamp is empty")
	}
	if !containsChar(pos.Timestamp, 'T') {
		t.Error("timestamp does not contain 'T'")
	}
}

func containsChar(s string, c rune) bool {
	for _, r := range s {
		if r == c {
			return true
		}
	}
	return false
}

package orbit

import "testing"

func TestPositionLatRange(t *testing.T) {
	pos, err := PropagatePosition(testTLEName, testTLELine1, testTLELine2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pos.Lat < -90 || pos.Lat > 90 {
		t.Errorf("lat %f out of range [-90, 90]", pos.Lat)
	}
}

func TestPositionLonRange(t *testing.T) {
	pos, err := PropagatePosition(testTLEName, testTLELine1, testTLELine2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pos.Lon < -180 || pos.Lon > 180 {
		t.Errorf("lon %f out of range [-180, 180]", pos.Lon)
	}
}

func TestPositionAltitudeRange(t *testing.T) {
	pos, err := PropagatePosition(testTLEName, testTLELine1, testTLELine2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pos.AltitudeKm < 200 || pos.AltitudeKm > 500 {
		t.Errorf("altitude %f out of range [200, 500]", pos.AltitudeKm)
	}
}

func TestPositionVelocityRange(t *testing.T) {
	pos, err := PropagatePosition(testTLEName, testTLELine1, testTLELine2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pos.VelocityKmh < 25000 || pos.VelocityKmh > 30000 {
		t.Errorf("velocity %f out of range [25000, 30000]", pos.VelocityKmh)
	}
}

func TestPositionTimestamp(t *testing.T) {
	pos, err := PropagatePosition(testTLEName, testTLELine1, testTLELine2)
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

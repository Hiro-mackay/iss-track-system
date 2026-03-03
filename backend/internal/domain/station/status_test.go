package station

import (
	"testing"

	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/tracking"
)

func TestComputeStatus_Period(t *testing.T) {
	status, err := ComputeStatus(testTLE(), 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	period := status.OrbitalParams.PeriodMin
	if period < 90 || period > 94 {
		t.Errorf("should have period ~92 min: got %.2f", period)
	}
}

func TestComputeStatus_Apogee(t *testing.T) {
	status, err := ComputeStatus(testTLE(), 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	apogee := status.OrbitalParams.ApogeeKm
	if apogee < 400 || apogee > 430 {
		t.Errorf("should have apogee 400-430 km: got %.2f", apogee)
	}
}

func TestComputeStatus_Perigee(t *testing.T) {
	status, err := ComputeStatus(testTLE(), 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	perigee := status.OrbitalParams.PerigeeKm
	if perigee < 400 || perigee > 430 {
		t.Errorf("should have perigee 400-430 km: got %.2f", perigee)
	}
}

func TestComputeStatus_Inclination(t *testing.T) {
	status, err := ComputeStatus(testTLE(), 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	inc := status.OrbitalParams.InclinationDeg
	if inc < 51.0 || inc > 52.0 {
		t.Errorf("should have inclination ~51.6 deg: got %.3f", inc)
	}
}

func TestComputeStatus_OrbitCount(t *testing.T) {
	status, err := ComputeStatus(testTLE(), 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.OrbitCount <= 40000 {
		t.Errorf("should have orbit count > 40000: got %d", status.OrbitCount)
	}
}

func TestComputeStatus_CrewCount(t *testing.T) {
	status, err := ComputeStatus(testTLE(), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.CrewCount != 7 {
		t.Errorf("should pass through crew count: got %d, want 7", status.CrewCount)
	}
}

func TestComputeStatus_LaunchYear(t *testing.T) {
	status, err := ComputeStatus(testTLE(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.LaunchYear != 1998 {
		t.Errorf("should have launch year 1998: got %d", status.LaunchYear)
	}
}

func TestComputeStatus_ZeroTLE(t *testing.T) {
	var tle tracking.TLE
	_, err := ComputeStatus(tle, 0)
	if err == nil {
		t.Error("should return error for zero-value TLE")
	}
}

package main

import (
	"testing"
)

const (
	testTLEName  = "ISS (ZARYA)"
	testTLELine1 = "1 25544U 98067A   24100.50000000  .00016717  00000-0  10270-3 0  9009"
	testTLELine2 = "2 25544  51.6400 200.0000 0001234  90.0000 270.0000 15.49000000400001"
)

func TestComputeISSStatus_Period(t *testing.T) {
	status, err := ComputeISSStatus(testTLEName, testTLELine1, testTLELine2, 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	period := status.OrbitalParams.PeriodMin
	if period < 90 || period > 94 {
		t.Errorf("should have period ~92 min: got %.2f", period)
	}
}

func TestComputeISSStatus_Apogee(t *testing.T) {
	status, err := ComputeISSStatus(testTLEName, testTLELine1, testTLELine2, 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	apogee := status.OrbitalParams.ApogeeKm
	if apogee < 400 || apogee > 430 {
		t.Errorf("should have apogee 400-430 km: got %.2f", apogee)
	}
}

func TestComputeISSStatus_Perigee(t *testing.T) {
	status, err := ComputeISSStatus(testTLEName, testTLELine1, testTLELine2, 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	perigee := status.OrbitalParams.PerigeeKm
	if perigee < 400 || perigee > 430 {
		t.Errorf("should have perigee 400-430 km: got %.2f", perigee)
	}
}

func TestComputeISSStatus_Inclination(t *testing.T) {
	status, err := ComputeISSStatus(testTLEName, testTLELine1, testTLELine2, 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	inc := status.OrbitalParams.InclinationDeg
	if inc < 51.0 || inc > 52.0 {
		t.Errorf("should have inclination ~51.6 deg: got %.3f", inc)
	}
}

func TestComputeISSStatus_OrbitCount(t *testing.T) {
	status, err := ComputeISSStatus(testTLEName, testTLELine1, testTLELine2, 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// RevolutionNumber in TLE is 40000, orbit count should be greater
	if status.OrbitCount <= 40000 {
		t.Errorf("should have orbit count > 40000: got %d", status.OrbitCount)
	}
}

func TestComputeISSStatus_CrewCount(t *testing.T) {
	status, err := ComputeISSStatus(testTLEName, testTLELine1, testTLELine2, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.CrewCount != 7 {
		t.Errorf("should pass through crew count: got %d, want 7", status.CrewCount)
	}
}

func TestComputeISSStatus_LaunchYear(t *testing.T) {
	status, err := ComputeISSStatus(testTLEName, testTLELine1, testTLELine2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.LaunchYear != 1998 {
		t.Errorf("should have launch year 1998: got %d", status.LaunchYear)
	}
}

func TestComputeISSStatus_InvalidTLE(t *testing.T) {
	_, err := ComputeISSStatus("bad", "bad", "bad", 0)
	if err == nil {
		t.Error("should return error for invalid TLE")
	}
}

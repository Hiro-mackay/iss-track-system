package pass

import "github.com/Hiro-mackay/iss-track-system/backend/internal/domain/tracking"

const (
	testTLEName  = "ISS (ZARYA)"
	testTLELine1 = "1 25544U 98067A   24100.50000000  .00016717  00000-0  10270-3 0  9009"
	testTLELine2 = "2 25544  51.6400 200.0000 0001234  90.0000 270.0000 15.49000000400001"
)

func testTLE() tracking.TLE {
	tle, err := tracking.NewTLE(testTLEName, testTLELine1, testTLELine2)
	if err != nil {
		panic("test TLE data is invalid: " + err.Error())
	}
	return tle
}

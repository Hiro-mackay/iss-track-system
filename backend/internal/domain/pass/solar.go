package pass

import (
	"math"
	"time"

	"github.com/akhenakh/sgp4"
)

const (
	earthRadiusKm = 6371.0
	deg2rad       = math.Pi / 180.0
	rad2deg       = 180.0 / math.Pi
)

// SunAltitude returns the sun's altitude angle in degrees for a given
// observer location and time. Uses a simplified solar position algorithm
// with approximately 1-degree accuracy.
func SunAltitude(lat, lon float64, t time.Time) float64 {
	jd := julianDay(t)
	n := jd - 2451545.0

	// Mean longitude and mean anomaly of the sun
	meanLon := math.Mod(280.460+0.9856474*n, 360.0)
	meanAnom := math.Mod(357.528+0.9856003*n, 360.0)

	// Ecliptic longitude
	gRad := meanAnom * deg2rad
	eclipLon := meanLon + 1.915*math.Sin(gRad) + 0.020*math.Sin(2*gRad)

	// Obliquity of the ecliptic
	obliquity := 23.439 - 0.0000004*n

	eclipRad := eclipLon * deg2rad
	obliqRad := obliquity * deg2rad

	// Right ascension and declination
	ra := math.Atan2(math.Cos(obliqRad)*math.Sin(eclipRad), math.Cos(eclipRad)) * rad2deg
	dec := math.Asin(math.Sin(obliqRad)*math.Sin(eclipRad)) * rad2deg

	// Greenwich Mean Sidereal Time
	gmst := math.Mod(280.46061837+360.98564736629*n, 360.0)

	// Hour angle
	ha := (gmst + lon - ra) * deg2rad

	latRad := lat * deg2rad
	decRad := dec * deg2rad

	alt := math.Asin(
		math.Sin(latRad)*math.Sin(decRad) +
			math.Cos(latRad)*math.Cos(decRad)*math.Cos(ha),
	)

	return alt * rad2deg
}

// IsPassVisible determines if an ISS pass is visible to an observer.
// A pass is visible when the satellite is sunlit AND the observer is in
// civil twilight or darker (sun altitude < -6 degrees).
func IsPassVisible(tle *sgp4.TLE, obsLat, obsLon, satAltKm float64, t time.Time) bool {
	// Check observer is in darkness (civil twilight)
	obsSunAlt := SunAltitude(obsLat, obsLon, t)
	if obsSunAlt >= -6.0 {
		return false
	}

	// Propagate TLE to get sub-satellite point
	eci, err := tle.FindPositionAtTime(t)
	if err != nil {
		return false
	}
	satLat, satLon, _ := eci.ToGeodetic()

	// Check satellite is sunlit: sun altitude at sub-satellite point
	// must be above the earth shadow angle
	satSunAlt := SunAltitude(satLat, satLon, t)
	shadowAngle := -math.Acos(earthRadiusKm/(earthRadiusKm+satAltKm)) * rad2deg

	return satSunAlt > shadowAngle
}

// julianDay converts a time.Time to Julian Day number.
func julianDay(t time.Time) float64 {
	y := float64(t.Year())
	m := float64(t.Month())
	d := float64(t.Day())
	h := float64(t.Hour()) + float64(t.Minute())/60.0 + float64(t.Second())/3600.0

	if m <= 2 {
		y--
		m += 12
	}

	a := math.Floor(y / 100.0)
	b := 2 - a + math.Floor(a/4.0)

	return math.Floor(365.25*(y+4716)) +
		math.Floor(30.6001*(m+1)) +
		d + h/24.0 + b - 1524.5
}

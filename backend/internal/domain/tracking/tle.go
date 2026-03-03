package tracking

import (
	"fmt"

	"github.com/akhenakh/sgp4"
)

// TLE is a Value Object representing parsed Two-Line Element data.
type TLE struct {
	name   string
	parsed *sgp4.TLE
}

// NewTLE creates a TLE Value Object by parsing the raw TLE lines.
// Returns an error if the TLE data cannot be parsed.
func NewTLE(name, line1, line2 string) (TLE, error) {
	tleStr := name + "\n" + line1 + "\n" + line2
	parsed, err := sgp4.ParseTLE(tleStr)
	if err != nil {
		return TLE{}, fmt.Errorf("parsing TLE: %w", err)
	}
	return TLE{name: name, parsed: parsed}, nil
}

// Name returns the satellite name from the TLE data.
func (t TLE) Name() string {
	return t.name
}

// Parsed returns the underlying sgp4.TLE for orbital computations.
func (t TLE) Parsed() *sgp4.TLE {
	return t.parsed
}

// IsValid reports whether the TLE contains parsed data (i.e. is not a zero value).
func (t TLE) IsValid() bool {
	return t.parsed != nil
}

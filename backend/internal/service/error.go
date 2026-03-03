package service

// TLEError wraps a TLE provider error so handlers can distinguish
// TLE unavailability (503) from computation errors (500).
type TLEError struct {
	Err error
}

// Error implements the error interface.
func (e *TLEError) Error() string {
	return e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *TLEError) Unwrap() error {
	return e.Err
}

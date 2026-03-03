package tracking

// TLEProvider provides access to cached TLE (Two-Line Element) data.
type TLEProvider interface {
	Get() (TLE, error)
	HasData() bool
	Refresh() error
}

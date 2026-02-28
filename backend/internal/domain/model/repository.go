package model

// TLEProvider provides access to cached TLE (Two-Line Element) data.
type TLEProvider interface {
	Get() (name, line1, line2 string, err error)
	HasData() bool
	Refresh() error
}

// CrewProvider provides access to cached ISS crew data.
type CrewProvider interface {
	Get() ([]CrewMember, error)
	Refresh() error
}

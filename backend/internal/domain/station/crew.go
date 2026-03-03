package station

// CrewMember represents a person currently aboard a spacecraft.
type CrewMember struct {
	Name  string `json:"name"`
	Craft string `json:"craft"`
}

// CrewProvider provides access to cached ISS crew data.
type CrewProvider interface {
	Get() ([]CrewMember, error)
	Refresh() error
}

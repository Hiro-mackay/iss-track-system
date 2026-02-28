package model

// CrewMember represents a person currently aboard a spacecraft.
type CrewMember struct {
	Name  string `json:"name"`
	Craft string `json:"craft"`
}

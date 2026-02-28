package query

import "github.com/Hiro-mackay/iss-track-system/backend/internal/domain/model"

// CrewQueryService provides crew data queries.
type CrewQueryService struct {
	crew model.CrewProvider
}

// NewCrewQueryService creates a new CrewQueryService.
func NewCrewQueryService(crew model.CrewProvider) *CrewQueryService {
	return &CrewQueryService{crew: crew}
}

// List returns the current ISS crew members.
func (s *CrewQueryService) List() ([]model.CrewMember, error) {
	return s.crew.Get()
}

package station

import (
	domain "github.com/Hiro-mackay/iss-track-system/backend/internal/domain/station"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/tracking"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/service"
)

// QueryService provides ISS crew and station status queries.
type QueryService struct {
	tle  tracking.TLEProvider
	crew domain.CrewProvider
}

// NewQueryService creates a new QueryService.
func NewQueryService(tle tracking.TLEProvider, crew domain.CrewProvider) *QueryService {
	return &QueryService{tle: tle, crew: crew}
}

// List returns the current ISS crew members.
func (s *QueryService) List() ([]domain.CrewMember, error) {
	return s.crew.Get()
}

// Status returns the current ISS station status.
func (s *QueryService) Status() (*domain.ISSStatus, error) {
	tle, err := s.tle.Get()
	if err != nil {
		return nil, &service.TLEError{Err: err}
	}

	crewCount := 0
	if members, err := s.crew.Get(); err == nil {
		crewCount = len(members)
	}

	status, err := domain.ComputeStatus(tle, crewCount)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

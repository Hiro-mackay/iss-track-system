package query

import (
	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/model"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/orbit"
)

// StationQueryService provides ISS station status queries.
type StationQueryService struct {
	tle  model.TLEProvider
	crew model.CrewProvider
}

// NewStationQueryService creates a new StationQueryService.
func NewStationQueryService(tle model.TLEProvider, crew model.CrewProvider) *StationQueryService {
	return &StationQueryService{tle: tle, crew: crew}
}

// Status returns the current ISS station status.
func (s *StationQueryService) Status() (*model.ISSStatus, error) {
	name, l1, l2, err := s.tle.Get()
	if err != nil {
		return nil, &TLEError{Err: err}
	}

	crewCount := 0
	if members, err := s.crew.Get(); err == nil {
		crewCount = len(members)
	}

	status, err := orbit.ComputeStatus(name, l1, l2, crewCount)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

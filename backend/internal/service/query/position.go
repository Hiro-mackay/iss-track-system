package query

import (
	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/model"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/orbit"
)

// PositionQueryService provides position and orbit track queries.
type PositionQueryService struct {
	tle model.TLEProvider
}

// NewPositionQueryService creates a new PositionQueryService.
func NewPositionQueryService(tle model.TLEProvider) *PositionQueryService {
	return &PositionQueryService{tle: tle}
}

// CurrentPosition returns the current ISS position.
func (s *PositionQueryService) CurrentPosition() (*model.ISSPosition, error) {
	name, l1, l2, err := s.tle.Get()
	if err != nil {
		return nil, &TLEError{Err: err}
	}

	pos, err := orbit.PropagatePosition(name, l1, l2)
	if err != nil {
		return nil, err
	}

	return &pos, nil
}

// OrbitTrack returns the ISS orbit track for the given number of minutes.
func (s *PositionQueryService) OrbitTrack(minutes int) ([]model.OrbitPoint, error) {
	name, l1, l2, err := s.tle.Get()
	if err != nil {
		return nil, &TLEError{Err: err}
	}

	return orbit.PropagateOrbitTrack(name, l1, l2, minutes)
}

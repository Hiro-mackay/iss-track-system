package tracking

import (
	domain "github.com/Hiro-mackay/iss-track-system/backend/internal/domain/tracking"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/service"
)

// QueryService provides position and orbit track queries.
type QueryService struct {
	tle domain.TLEProvider
}

// NewQueryService creates a new QueryService.
func NewQueryService(tle domain.TLEProvider) *QueryService {
	return &QueryService{tle: tle}
}

// CurrentPosition returns the current ISS position.
func (s *QueryService) CurrentPosition() (*domain.ISSPosition, error) {
	tle, err := s.tle.Get()
	if err != nil {
		return nil, &service.TLEError{Err: err}
	}

	pos, err := tle.PropagatePosition()
	if err != nil {
		return nil, err
	}

	return &pos, nil
}

// OrbitTrack returns the ISS orbit track for the given number of minutes.
func (s *QueryService) OrbitTrack(minutes int) ([]domain.OrbitPoint, error) {
	tle, err := s.tle.Get()
	if err != nil {
		return nil, &service.TLEError{Err: err}
	}

	return tle.PropagateOrbitTrack(minutes)
}

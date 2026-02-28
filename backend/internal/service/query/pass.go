package query

import (
	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/model"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/orbit"
)

// PassQueryService provides pass prediction queries.
type PassQueryService struct {
	tle model.TLEProvider
}

// NewPassQueryService creates a new PassQueryService.
func NewPassQueryService(tle model.TLEProvider) *PassQueryService {
	return &PassQueryService{tle: tle}
}

// Predict returns ISS pass predictions for the given observer location and time range.
func (s *PassQueryService) Predict(lat, lon float64, days int) ([]model.PassPrediction, error) {
	name, l1, l2, err := s.tle.Get()
	if err != nil {
		return nil, &TLEError{Err: err}
	}

	return orbit.PredictPasses(name, l1, l2, lat, lon, days)
}

package pass

import (
	domain "github.com/Hiro-mackay/iss-track-system/backend/internal/domain/pass"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/domain/tracking"
	"github.com/Hiro-mackay/iss-track-system/backend/internal/service"
)

// QueryService provides pass prediction queries.
type QueryService struct {
	tle tracking.TLEProvider
}

// NewQueryService creates a new QueryService.
func NewQueryService(tle tracking.TLEProvider) *QueryService {
	return &QueryService{tle: tle}
}

// Predict returns ISS pass predictions for the given observer location and time range.
func (s *QueryService) Predict(lat, lon float64, days int) ([]domain.PassPrediction, error) {
	tle, err := s.tle.Get()
	if err != nil {
		return nil, &service.TLEError{Err: err}
	}

	return domain.PredictPasses(tle, lat, lon, days)
}

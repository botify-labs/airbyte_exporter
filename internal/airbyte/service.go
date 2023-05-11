// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package airbyte

// Service handles domain operations for gathering metrics from Airbyte's PostgreSQL database.
type Service struct {
	r *Repository
}

// NewService initializes and returns an Airbyte Service.
func NewService(r *Repository) *Service {
	return &Service{
		r: r,
	}
}

// GatherMetrics gathers and returns metrics from Airbyte's PostgreSQL database.
func (s *Service) GatherMetrics() (*Metrics, error) {
	jobPending, err := s.r.JobPending()
	if err != nil {
		return &Metrics{}, err
	}

	jobRunning, err := s.r.JobRunning()
	if err != nil {
		return &Metrics{}, err
	}
	jobRunningOrphan, err := s.r.JobRunningOrphan()
	if err != nil {
		return &Metrics{}, err
	}

	return &Metrics{
		JobPending:       jobPending,
		JobRunning:       jobRunning,
		JobRunningOrphan: jobRunningOrphan,
	}, nil
}

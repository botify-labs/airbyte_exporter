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
	connections, err := s.r.ConnectionsCount()
	if err != nil {
		return &Metrics{}, err
	}

	jobsCompleted, err := s.r.JobsCompletedCount()
	if err != nil {
		return &Metrics{}, err
	}

	jobsPending, err := s.r.JobsPendingCount()
	if err != nil {
		return &Metrics{}, err
	}

	jobsRunning, err := s.r.JobsRunningCount()
	if err != nil {
		return &Metrics{}, err
	}

	return &Metrics{
		Connections:   connections,
		JobsCompleted: jobsCompleted,
		JobsPending:   jobsPending,
		JobsRunning:   jobsRunning,
	}, nil
}

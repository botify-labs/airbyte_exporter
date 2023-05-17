// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package airbyte

import "github.com/jmoiron/sqlx"

// Repository provides an abstraction layer to perform SQL queries against the
// Airbyte PostgreSQL database.
type Repository struct {
	db *sqlx.DB
}

// NewRepository initializes and returns an Airbyte Repository.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// countBySourceAndStatusQuery provides a helper to run a SQL query that returns rows to be marshaled
// as a slice of CountBySourceAndStatus.
func (r *Repository) countBySourceAndStatusQuery(query string) ([]CountBySourceAndStatus, error) {
	rows, err := r.db.Queryx(query)
	if err != nil {
		return []CountBySourceAndStatus{}, err
	}

	var jobStatuses []CountBySourceAndStatus

	for rows.Next() {
		var jobStatus CountBySourceAndStatus

		if err := rows.StructScan(&jobStatus); err != nil {
			return []CountBySourceAndStatus{}, err
		}

		jobStatuses = append(jobStatuses, jobStatus)
	}

	return jobStatuses, nil
}

// ConnectionsBySourceName returns the count of Airbyte connections, grouped by  source and status.
func (r *Repository) ConnectionsBySourceName() ([]CountBySourceAndStatus, error) {
	query := `
	SELECT ad.name as source, c.status, COUNT(c.status)
	FROM connection c
	LEFT JOIN actor a ON c.source_id = a.id
	LEFT JOIN actor_definition ad ON a.actor_definition_id = ad.id
	GROUP BY ad.name, c.status
	ORDER BY ad.name, c.status
	`

	return r.countBySourceAndStatusQuery(query)
}

// JobsCompletedBySourceName returns the count of completed Airbyte jobs, grouped by source and status.
func (r *Repository) JobsCompletedBySourceName() ([]CountBySourceAndStatus, error) {
	query := `
	SELECT ad.name as source, j.status, COUNT(j.status)
	FROM jobs j
	LEFT JOIN connection c ON j.scope = CAST(c.id AS VARCHAR(255))
	LEFT JOIN actor a ON c.source_id = a.id
	LEFT JOIN actor_definition ad ON a.actor_definition_id = ad.id
	WHERE j.status IN ('cancelled', 'failed', 'succeeded')
	GROUP BY ad.name, j.status
	ORDER BY ad.name, j.status
	`

	return r.countBySourceAndStatusQuery(query)
}

// JobsPendingBySourceName returns the count of pending Airbyte jobs, grouped by source and status.
func (r *Repository) JobsPendingBySourceName() ([]CountBySourceAndStatus, error) {
	query := `
	SELECT ad.name as source, j.status, COUNT(j.status)
	FROM jobs j
	LEFT JOIN connection c ON CAST(c.id AS VARCHAR(255)) = j.scope
	LEFT JOIN actor a ON c.source_id = a.id
	LEFT JOIN actor_definition ad ON a.actor_definition_id = ad.id
	WHERE j.status = 'pending'
	GROUP BY ad.name, j.status
	ORDER BY ad.name, j.status
	`

	return r.countBySourceAndStatusQuery(query)
}

// JobsRunningBySourceName returns the count of running Airbyte jobs, grouped by source and status.
func (r *Repository) JobsRunningBySourceName() ([]CountBySourceAndStatus, error) {
	query := `
	SELECT ad.name as source, j.status, COUNT(j.status)
	FROM jobs j
	LEFT JOIN attempts att ON att.job_id = j.id
	LEFT JOIN connection c ON j.scope = CAST(c.id AS VARCHAR(255))
	LEFT JOIN actor a ON c.source_id = a.id
	LEFT JOIN actor_definition ad ON a.actor_definition_id = ad.id
	WHERE j.status = 'running'
	AND   att.status = 'running'
	GROUP BY ad.name, j.status
	ORDER BY ad.name, j.status
	`

	return r.countBySourceAndStatusQuery(query)
}

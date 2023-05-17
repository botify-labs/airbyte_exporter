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

// jobStatusCountQuery provides a helper to run a SQL query that returns rows to be marshaled
// as a slice of JobStatusCount.
func (r *Repository) jobStatusCountQuery(query string) ([]JobStatusCount, error) {
	rows, err := r.db.Queryx(query)
	if err != nil {
		return []JobStatusCount{}, err
	}

	var jobStatuses []JobStatusCount

	for rows.Next() {
		var jobStatus JobStatusCount

		if err := rows.StructScan(&jobStatus); err != nil {
			return []JobStatusCount{}, err
		}

		jobStatuses = append(jobStatuses, jobStatus)
	}

	return jobStatuses, nil
}

// JobsPendingBySourceName returns the count of pending Airbyte jobs, grouped by source and status.
func (r *Repository) JobsPendingBySourceName() ([]JobStatusCount, error) {
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

	return r.jobStatusCountQuery(query)
}

// JobsRunningBySourceName returns the count of running Airbyte jobs, grouped by source and status.
func (r *Repository) JobsRunningBySourceName() ([]JobStatusCount, error) {
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

	return r.jobStatusCountQuery(query)
}

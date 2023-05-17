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

// countQuery provides a helper to run a SQL query that returns rows to be marshaled
// as a slice of ConnectionCount.
func (r *Repository) connectionCountQuery(query string) ([]ConnectionCount, error) {
	rows, err := r.db.Queryx(query)
	if err != nil {
		return []ConnectionCount{}, err
	}

	var connectionCounts []ConnectionCount

	for rows.Next() {
		var connectionCount ConnectionCount

		if err := rows.StructScan(&connectionCount); err != nil {
			return []ConnectionCount{}, err
		}

		connectionCounts = append(connectionCounts, connectionCount)
	}

	return connectionCounts, nil
}

// jobCountQuery provides a helper to run a SQL query that returns rows to be marshaled
// as a slice of JobCount.
func (r *Repository) jobCountQuery(query string) ([]JobCount, error) {
	rows, err := r.db.Queryx(query)
	if err != nil {
		return []JobCount{}, err
	}

	var jobCounts []JobCount

	for rows.Next() {
		var jobCount JobCount

		if err := rows.StructScan(&jobCount); err != nil {
			return []JobCount{}, err
		}

		jobCounts = append(jobCounts, jobCount)
	}

	return jobCounts, nil
}

// ConnectionsCount returns the count of Airbyte connections, grouped by destination, source and status.
func (r *Repository) ConnectionsCount() ([]ConnectionCount, error) {
	query := `
	SELECT ad1.name as destination, ad2.name as source, c.status, COUNT(c.status)
	FROM connection c
	LEFT JOIN actor a1 ON c.destination_id = a1.id
	LEFT JOIN actor_definition ad1 ON a1.actor_definition_id = ad1.id
	LEFT JOIN actor a2 ON c.source_id = a2.id
	LEFT JOIN actor_definition ad2 ON a2.actor_definition_id = ad2.id
	GROUP BY ad1.name, ad2.name, c.status
	ORDER BY ad1.name, ad2.name, c.status
	`

	return r.connectionCountQuery(query)
}

// JobsCompletedCount returns the count of completed Airbyte jobs, grouped by destination, source, type and status.
func (r *Repository) JobsCompletedCount() ([]JobCount, error) {
	query := `
	SELECT ad1.name as destination, ad2.name as source, j.config_type, j.status, COUNT(j.status)
	FROM jobs j
	LEFT JOIN connection c ON j.scope = CAST(c.id AS VARCHAR(255))
	LEFT JOIN actor a1 ON c.destination_id = a1.id
	LEFT JOIN actor_definition ad1 ON a1.actor_definition_id = ad1.id
	LEFT JOIN actor a2 ON c.source_id = a2.id
	LEFT JOIN actor_definition ad2 ON a2.actor_definition_id = ad2.id
	WHERE j.status IN ('cancelled', 'failed', 'succeeded')
	GROUP BY ad1.name, ad2.name, j.config_type, j.status
	ORDER BY ad1.name, ad2.name, j.config_type, j.status
	`

	return r.jobCountQuery(query)
}

// JobsPendingCount returns the count of pending Airbyte jobs, grouped by destination, source and type.
func (r *Repository) JobsPendingCount() ([]JobCount, error) {
	query := `
	SELECT ad1.name as destination, ad2.name as source, j.config_type, j.status, COUNT(j.status)
	FROM jobs j
	LEFT JOIN connection c ON CAST(c.id AS VARCHAR(255)) = j.scope
	LEFT JOIN actor a1 ON c.destination_id = a1.id
	LEFT JOIN actor_definition ad1 ON a1.actor_definition_id = ad1.id
	LEFT JOIN actor a2 ON c.source_id = a2.id
	LEFT JOIN actor_definition ad2 ON a2.actor_definition_id = ad2.id
	WHERE j.status = 'pending'
	GROUP BY ad1.name, ad2.name, j.config_type, j.status
	ORDER BY ad1.name, ad2.name, j.config_type, j.status
	`

	return r.jobCountQuery(query)
}

// JobsRunningCount returns the count of running Airbyte jobs, grouped by destination, source and type.
func (r *Repository) JobsRunningCount() ([]JobCount, error) {
	query := `
	SELECT ad1.name as destination, ad2.name as source, j.config_type, j.status, COUNT(j.status)
	FROM jobs j
	LEFT JOIN attempts att ON att.job_id = j.id
	LEFT JOIN connection c ON j.scope = CAST(c.id AS VARCHAR(255))
	LEFT JOIN actor a1 ON c.destination_id = a1.id
	LEFT JOIN actor_definition ad1 ON a1.actor_definition_id = ad1.id
	LEFT JOIN actor a2 ON c.source_id = a2.id
	LEFT JOIN actor_definition ad2 ON a2.actor_definition_id = ad2.id
	WHERE j.status = 'running'
	AND   att.status = 'running'
	GROUP BY ad1.name, ad2.name, j.config_type, j.status
	ORDER BY ad1.name, ad2.name, j.config_type, j.status
	`

	return r.jobCountQuery(query)
}

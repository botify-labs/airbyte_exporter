// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package airbyte

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides an abstraction layer to perform SQL queries against the
// Airbyte PostgreSQL database.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository initializes and returns an Airbyte Repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// actorCountQuery provides a helper to run a SQL query that returns rows to be marshaled
// as a slice of ActorCount.
func (r *Repository) actorCountQuery(query string) ([]ActorCount, error) {
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return []ActorCount{}, err
	}

	var actorCounts []ActorCount
	if err := pgxscan.ScanAll(&actorCounts, rows); err != nil {
		return []ActorCount{}, err
	}

	return actorCounts, nil
}

// connectionCountQuery provides a helper to run a SQL query that returns rows to be marshaled
// as a slice of ConnectionCount.
func (r *Repository) connectionCountQuery(query string) ([]ConnectionCount, error) {
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return []ConnectionCount{}, err
	}

	var connectionCounts []ConnectionCount
	if err := pgxscan.ScanAll(&connectionCounts, rows); err != nil {
		return []ConnectionCount{}, err
	}

	return connectionCounts, nil
}

// connectionSyncAgeQuery provides a helper to run a SQL query that returns rows to be marshaled
// as a slice of ConnectionSyncAge.
func (r *Repository) connectionSyncAgeQuery(query string) ([]ConnectionSyncAge, error) {
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return []ConnectionSyncAge{}, err
	}

	var connectionSyncAges []ConnectionSyncAge
	if err := pgxscan.ScanAll(&connectionSyncAges, rows); err != nil {
		return []ConnectionSyncAge{}, err
	}

	return connectionSyncAges, nil
}

// jobCountQuery provides a helper to run a SQL query that returns rows to be marshaled
// as a slice of JobCount.
func (r *Repository) jobCountQuery(query string) ([]JobCount, error) {
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return []JobCount{}, err
	}

	var jobCounts []JobCount
	if err := pgxscan.ScanAll(&jobCounts, rows); err != nil {
		return []JobCount{}, err
	}

	return jobCounts, nil
}

// ConnectionsCount returns the count of Airbyte connections, grouped by destination, source and status.
func (r *Repository) ConnectionsCount() ([]ConnectionCount, error) {
	query := `
	SELECT ad1.name as destination, ad2.name as source, c.status, COUNT(c.status)
	FROM connection c
	JOIN actor a1 ON c.destination_id = a1.id
	JOIN actor_definition ad1 ON a1.actor_definition_id = ad1.id
	JOIN actor a2 ON c.source_id = a2.id
	JOIN actor_definition ad2 ON a2.actor_definition_id = ad2.id
	GROUP BY ad1.name, ad2.name, c.status
	ORDER BY ad1.name, ad2.name, c.status
	`

	return r.connectionCountQuery(query)
}

// ConnectionsLastSuccessfulSyncAge returns the age of the last successful sync job attempt
// for active connections.
func (r *Repository) ConnectionsLastSuccessfulSyncAge() ([]ConnectionSyncAge, error) {
	query := `
	WITH j AS (
		SELECT scope, max(updated_at) AS updated_at
		FROM  jobs
		WHERE config_type = 'sync'
		AND   status = 'succeeded'
		GROUP BY scope
	)
	SELECT c.id, COALESCE(c.schedule_type, 'manual') AS connection_schedule_type, ad1.name as destination, ad2.name as source, EXTRACT(EPOCH FROM AGE(NOW(), j.updated_at))/3600 as hours
	FROM connection c
	JOIN j ON j.scope = CAST(c.id AS VARCHAR(255))
	JOIN actor a1 ON c.destination_id = a1.id
	JOIN actor_definition ad1 ON a1.actor_definition_id = ad1.id
	JOIN actor a2 ON c.source_id = a2.id
	JOIN actor_definition ad2 ON a2.actor_definition_id = ad2.id
	WHERE c.status = 'active'
	`

	return r.connectionSyncAgeQuery(query)
}

// SourcesCount returns the count of Airbyte sources, grouped by actor connector and status.
func (r *Repository) SourcesCount() ([]ActorCount, error) {
	query := `
	SELECT ad.name as actor, a.tombstone, COUNT(a.tombstone)
	FROM actor a
	JOIN actor_definition ad ON a.actor_definition_id = ad.id
	WHERE a.actor_type = 'source'
	GROUP BY ad.name, a.tombstone
	ORDER BY ad.name, a.tombstone
	`
	return r.actorCountQuery(query)
}

// DestinationsCount returns the count of Airbyte sources, grouped by actor connector and status.
func (r *Repository) DestinationsCount() ([]ActorCount, error) {
	query := `
	SELECT ad.name as actor, a.tombstone, COUNT(a.tombstone)
	FROM actor a
	JOIN actor_definition ad ON a.actor_definition_id = ad.id
	WHERE a.actor_type = 'destination'
	GROUP BY ad.name, a.tombstone
	ORDER BY ad.name, a.tombstone
	`
	return r.actorCountQuery(query)
}

// JobsCompletedCount returns the count of completed Airbyte jobs, grouped by destination, source, type and status.
func (r *Repository) JobsCompletedCount() ([]JobCount, error) {
	query := `
	SELECT ad1.name as destination, ad2.name as source, COALESCE(c.schedule_type, 'manual') AS connection_schedule_type, j.config_type, j.status, COUNT(j.status)
	FROM jobs j
	JOIN connection c ON j.scope = CAST(c.id AS VARCHAR(255))
	JOIN actor a1 ON c.destination_id = a1.id
	JOIN actor_definition ad1 ON a1.actor_definition_id = ad1.id
	JOIN actor a2 ON c.source_id = a2.id
	JOIN actor_definition ad2 ON a2.actor_definition_id = ad2.id
	WHERE j.status IN ('cancelled', 'failed', 'succeeded')
	GROUP BY ad1.name, ad2.name, connection_schedule_type, j.config_type, j.status
	ORDER BY ad1.name, ad2.name, connection_schedule_type, j.config_type, j.status
	`

	return r.jobCountQuery(query)
}

// JobsPendingCount returns the count of pending Airbyte jobs, grouped by destination, source and type.
func (r *Repository) JobsPendingCount() ([]JobCount, error) {
	query := `
	SELECT ad1.name as destination, ad2.name as source, COALESCE(c.schedule_type, 'manual') AS connection_schedule_type, j.config_type, j.status, COUNT(j.status)
	FROM jobs j
	JOIN connection c ON CAST(c.id AS VARCHAR(255)) = j.scope
	JOIN actor a1 ON c.destination_id = a1.id
	JOIN actor_definition ad1 ON a1.actor_definition_id = ad1.id
	JOIN actor a2 ON c.source_id = a2.id
	JOIN actor_definition ad2 ON a2.actor_definition_id = ad2.id
	WHERE j.status = 'pending'
	GROUP BY ad1.name, ad2.name, connection_schedule_type, j.config_type, j.status
	ORDER BY ad1.name, ad2.name, connection_schedule_type, j.config_type, j.status
	`

	return r.jobCountQuery(query)
}

// JobsRunningCount returns the count of running Airbyte jobs, grouped by destination, source and type.
func (r *Repository) JobsRunningCount() ([]JobCount, error) {
	query := `
	SELECT ad1.name as destination, ad2.name as source, COALESCE(c.schedule_type, 'manual') AS connection_schedule_type, j.config_type, j.status, COUNT(j.status)
	FROM jobs j
	JOIN attempts att ON att.job_id = j.id
	JOIN connection c ON j.scope = CAST(c.id AS VARCHAR(255))
	JOIN actor a1 ON c.destination_id = a1.id
	JOIN actor_definition ad1 ON a1.actor_definition_id = ad1.id
	JOIN actor a2 ON c.source_id = a2.id
	JOIN actor_definition ad2 ON a2.actor_definition_id = ad2.id
	WHERE j.status = 'running'
	AND   att.status = 'running'
	GROUP BY ad1.name, ad2.name, connection_schedule_type, j.config_type, j.status
	ORDER BY ad1.name, ad2.name, connection_schedule_type, j.config_type, j.status
	`

	return r.jobCountQuery(query)
}

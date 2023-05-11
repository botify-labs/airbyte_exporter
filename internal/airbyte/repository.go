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

func (r *Repository) JobPending() (count uint, err error) {
	query := `
	SELECT COUNT(*)
	FROM jobs
	JOIN connection
	ON   CAST(connection.id AS VARCHAR(255)) = jobs.scope
	WHERE jobs.status = 'pending'
	`

	err = r.db.Get(
		&count,
		query,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) JobRunning() (count uint, err error) {
	query := `
	SELECT COUNT(*)
	FROM jobs
	JOIN connection
	ON   CAST(connection.id AS VARCHAR(255)) = jobs.scope
	JOIN attempts
	ON   attempts.job_id = jobs.id
	WHERE jobs.status = 'running'
	AND   attempts.status = 'running'
	AND   connection.status = 'active'
	`

	err = r.db.Get(
		&count,
		query,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) JobRunningOrphan() (count uint, err error) {
	query := `
	SELECT COUNT(*)
	FROM jobs
	JOIN connection
	ON   CAST(connection.id AS VARCHAR(255)) = jobs.scope
	WHERE jobs.status = 'running'
	AND   connection.status != 'active'
	`

	err = r.db.Get(
		&count,
		query,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

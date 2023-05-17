// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package airbyte

// Metrics represents available Airbyte metrics.
type Metrics struct {
	// Airbyte connections
	Connections []CountBySourceAndStatus

	// Airbyte jobs
	JobsCompleted []CountBySourceAndStatus
	JobsPending   []CountBySourceAndStatus
	JobsRunning   []CountBySourceAndStatus
}

// CountBySourceAndStatus holds a count of Airbyte items, grouped by source name and status.
type CountBySourceAndStatus struct {
	Source string `db:"source"`
	Status string `db:"status"`
	Count  uint   `db:"count"`
}

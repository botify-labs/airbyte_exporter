// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package airbyte

// Metrics represents available Airbyte metrics.
type Metrics struct {
	// Airbyte jobs
	JobsPending []JobStatusCount
	JobsRunning []JobStatusCount
}

// JobStatusCount holds a count of Airbyte jobs, grouped by source name and status.
type JobStatusCount struct {
	Source string `db:"source"`
	Status string `db:"status"`
	Count  uint   `db:"count"`
}

// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package airbyte

// Metrics represents available Airbyte metrics.
type Metrics struct {
	// Airbyte connections
	Connections []ConnectionCount

	// Airbyte jobs
	JobsCompleted []JobCount
	JobsPending   []JobCount
	JobsRunning   []JobCount
}

// ConnectionCount holds a count of Airbyte connections, grouped by destination connector, source connector and status.
type ConnectionCount struct {
	DestinationConnector string `db:"destination"`
	SourceConnector      string `db:"source"`
	Status               string `db:"status"`
	Count                uint   `db:"count"`
}

// JobCount holds a count of Airbyte jobs, grouped by destination connector, source connector, type and status.
type JobCount struct {
	DestinationConnector string `db:"destination"`
	SourceConnector      string `db:"source"`
	Type                 string `db:"config_type"`
	Status               string `db:"status"`
	Count                uint   `db:"count"`
}

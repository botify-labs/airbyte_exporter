// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package airbyte

// Metrics represents available Airbyte metrics.
type Metrics struct {
	// Airbyte connections
	Connections []ItemCount

	// Airbyte jobs
	JobsCompleted []ItemCount
	JobsPending   []ItemCount
	JobsRunning   []ItemCount
}

// ItemCount holds a count of Airbyte items, grouped by destination connector, source connector and status.
type ItemCount struct {
	DestinationConnector string `db:"destination"`
	SourceConnector      string `db:"source"`
	Status               string `db:"status"`
	Count                uint   `db:"count"`
}

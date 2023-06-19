// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package airbyte

import (
	"fmt"
	"math"
	"time"
)

// Metrics represents available Airbyte metrics.
type Metrics struct {
	// Airbyte connections
	Connections                       []ConnectionCount
	ConnectionsLastSuccessfulSyncAges []ConnectionSyncAge

	// Airbyte connectors
	Sources      []ActorCount
	Destinations []ActorCount

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

// ConnectionSyncAge holds the age of the last job attempt for a single Airbyte Connection.
type ConnectionSyncAge struct {
	ID                   string  `db:"id"`
	DestinationConnector string  `db:"destination"`
	SourceConnector      string  `db:"source"`
	Hours                float64 `db:"hours"` // no Scanner for time.Duration, storing as a raw string
}

// Age returns the duration since the last job attempt.
func (csa *ConnectionSyncAge) Age() (time.Duration, error) {
	return time.ParseDuration(fmt.Sprintf("%fh", math.Round(csa.Hours)))
}

// ActorCount holds a count of Airbyte actors, grouped by actor connector and status.
type ActorCount struct {
	ActorConnector string `db:"actor"`
	Tombstone      bool   `db:"tombstone"`
	Count          uint   `db:"count"`
}

// JobCount holds a count of Airbyte jobs, grouped by destination connector, source connector, type and status.
type JobCount struct {
	DestinationConnector string `db:"destination"`
	SourceConnector      string `db:"source"`
	Type                 string `db:"config_type"`
	Status               string `db:"status"`
	Count                uint   `db:"count"`
}

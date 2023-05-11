// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package airbyte

// Metrics represents available Airbyte metrics.
type Metrics struct {
	// Airbyte jobs
	JobPending       uint
	JobRunning       uint
	JobRunningOrphan uint
}

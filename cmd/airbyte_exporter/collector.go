// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "airbyte"
)

// collector collects and exposes Airbyte metrics.
type collector struct {
	// Airbyte jobs
	jobPending        *prometheus.Desc
	jobRunning        *prometheus.Desc
	jobRunningOrphans *prometheus.Desc
}

// NewCollector initializes and returns a Prometheus collector for Airbyte metrics.
func NewCollector() *collector {
	return &collector{
		jobPending: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "job_pending"),
			"Pending jobs",
			nil,
			nil,
		),
		jobRunning: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "job_running"),
			"Running jobs",
			nil,
			nil,
		),
		jobRunningOrphans: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "job_running_orphan"),
			"Running jobs associated with an inactive or deprecated connection",
			nil,
			nil,
		),
	}
}

// Describe publishes the description of each Airbyte metric to a metrics
// channel.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.jobPending
	ch <- c.jobRunning
	ch <- c.jobRunningOrphans
}

// Collect gathers metrics from Airbyte.
func (c *collector) Collect(ch chan<- prometheus.Metric) {}

// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"

	"github.com/virtualtam/airbyte_exporter/internal/airbyte"
)

const (
	namespace = "airbyte"
)

// collector collects and exposes Airbyte metrics.
type collector struct {
	// Services
	airbyteService *airbyte.Service

	// Airbyte connections
	connections *prometheus.Desc
	// Airbyte jobs
	jobsCompleted *prometheus.Desc
	jobsPending   *prometheus.Desc
	jobsRunning   *prometheus.Desc
}

// NewCollector initializes and returns a Prometheus collector for Airbyte metrics.
func NewCollector(airbyteService *airbyte.Service) *collector {
	return &collector{
		airbyteService: airbyteService,

		connections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connections_total"),
			"Connections (total)",
			[]string{"source", "status"},
			nil,
		),

		jobsCompleted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "jobs_completed_total"),
			"Completed jobs (total)",
			[]string{"source", "status"},
			nil,
		),
		jobsPending: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "jobs_pending"),
			"Pending jobs",
			[]string{"source"},
			nil,
		),
		jobsRunning: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "jobs_running"),
			"Running jobs",
			[]string{"source"},
			nil,
		),
	}
}

// Describe publishes the description of each Airbyte metric to a metrics
// channel.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.connections
	ch <- c.jobsCompleted
	ch <- c.jobsPending
	ch <- c.jobsRunning
}

// Collect gathers metrics from Airbyte.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	metrics, err := c.airbyteService.GatherMetrics()
	if err != nil {
		log.Error().Err(err).Msg("failed to gather metrics")
	}

	// Counters
	for _, connections := range metrics.Connections {
		ch <- prometheus.MustNewConstMetric(
			c.connections,
			prometheus.CounterValue,
			float64(connections.Count),
			connections.Source,
			connections.Status,
		)
	}

	for _, jobsCompleted := range metrics.JobsCompleted {
		ch <- prometheus.MustNewConstMetric(
			c.jobsCompleted,
			prometheus.CounterValue,
			float64(jobsCompleted.Count),
			jobsCompleted.Source,
			jobsCompleted.Status,
		)
	}

	// Gauges
	for _, jobsPending := range metrics.JobsPending {
		ch <- prometheus.MustNewConstMetric(c.jobsPending, prometheus.GaugeValue, float64(jobsPending.Count), jobsPending.Source)
	}

	for _, jobsRunning := range metrics.JobsRunning {
		ch <- prometheus.MustNewConstMetric(c.jobsRunning, prometheus.GaugeValue, float64(jobsRunning.Count), jobsRunning.Source)
	}
}

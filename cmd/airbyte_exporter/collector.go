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

	// Airbyte jobs
	jobPending       *prometheus.Desc
	jobRunning       *prometheus.Desc
	jobRunningOrphan *prometheus.Desc
}

// NewCollector initializes and returns a Prometheus collector for Airbyte metrics.
func NewCollector(airbyteService *airbyte.Service) *collector {
	return &collector{
		airbyteService: airbyteService,

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
		jobRunningOrphan: prometheus.NewDesc(
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
	ch <- c.jobRunningOrphan
}

// Collect gathers metrics from Airbyte.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	metrics, err := c.airbyteService.GatherMetrics()
	if err != nil {
		log.Error().Err(err).Msg("failed to gather metrics")
	}

	// Gauges
	ch <- prometheus.MustNewConstMetric(c.jobPending, prometheus.GaugeValue, float64(metrics.JobPending))
	ch <- prometheus.MustNewConstMetric(c.jobRunning, prometheus.GaugeValue, float64(metrics.JobRunning))
	ch <- prometheus.MustNewConstMetric(c.jobRunningOrphan, prometheus.GaugeValue, float64(metrics.JobRunningOrphan))
}

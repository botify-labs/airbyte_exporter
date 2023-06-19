// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"

	"github.com/botify-labs/airbyte_exporter/v2/internal/airbyte"
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

	// Airbyte connectors
	sources      *prometheus.Desc
	destinations *prometheus.Desc

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
			prometheus.BuildFQName(namespace, "", "connections"),
			"Connections",
			[]string{"destination_connector", "source_connector", "status"},
			nil,
		),
		sources: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "sources"),
			"Sources",
			[]string{"source_connector", "tombstone"},
			nil,
		),
		destinations: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "destinations"),
			"Destinations",
			[]string{"destination_connector", "tombstone"},
			nil,
		),

		jobsCompleted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "jobs_completed_total"),
			"Completed jobs (total)",
			[]string{"destination_connector", "source_connector", "type", "status"},
			nil,
		),
		jobsPending: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "jobs_pending"),
			"Pending jobs",
			[]string{"destination_connector", "source_connector", "type"},
			nil,
		),
		jobsRunning: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "jobs_running"),
			"Running jobs",
			[]string{"destination_connector", "source_connector", "type"},
			nil,
		),
	}
}

// Describe publishes the description of each Airbyte metric to a metrics
// channel.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.connections
	ch <- c.sources
	ch <- c.destinations
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
	for _, jobsCompleted := range metrics.JobsCompleted {
		ch <- prometheus.MustNewConstMetric(
			c.jobsCompleted,
			prometheus.CounterValue,
			float64(jobsCompleted.Count),
			jobsCompleted.DestinationConnector,
			jobsCompleted.SourceConnector,
			jobsCompleted.Type,
			jobsCompleted.Status,
		)
	}

	// Gauges
	for _, connections := range metrics.Connections {
		ch <- prometheus.MustNewConstMetric(
			c.connections,
			prometheus.GaugeValue,
			float64(connections.Count),
			connections.DestinationConnector,
			connections.SourceConnector,
			connections.Status,
		)
	}

	for _, sources := range metrics.Sources {
		ch <- prometheus.MustNewConstMetric(
			c.sources,
			prometheus.GaugeValue,
			float64(sources.Count),
			sources.ActorConnector,
			strconv.FormatBool(sources.Tombstone),
		)
	}

	for _, destinations := range metrics.Destinations {
		ch <- prometheus.MustNewConstMetric(
			c.destinations,
			prometheus.GaugeValue,
			float64(destinations.Count),
			destinations.ActorConnector,
			strconv.FormatBool(destinations.Tombstone),
		)
	}

	for _, jobsPending := range metrics.JobsPending {
		ch <- prometheus.MustNewConstMetric(
			c.jobsPending,
			prometheus.GaugeValue,
			float64(jobsPending.Count),
			jobsPending.DestinationConnector,
			jobsPending.SourceConnector,
			jobsPending.Type,
		)
	}

	for _, jobsRunning := range metrics.JobsRunning {
		ch <- prometheus.MustNewConstMetric(
			c.jobsRunning,
			prometheus.GaugeValue,
			float64(jobsRunning.Count),
			jobsRunning.DestinationConnector,
			jobsRunning.SourceConnector,
			jobsRunning.Type,
		)
	}

	// Histograms
	connectionsLastSuccessfulSyncHistogramVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "connections_last_successful_sync_age_hours",
			Help:      "Age of the last successful sync job (hours)",
			Buckets:   []float64{6, 12, 18, 24, 48, 72, 168},
		},
		[]string{"destination_connector", "source_connector"},
	)

	for _, connectionLastSuccessfulSyncAge := range metrics.ConnectionsLastSuccessfulSyncAges {
		age, err := connectionLastSuccessfulSyncAge.Age()
		if err != nil {
			log.
				Error().
				Err(err).
				Str("connection_id", connectionLastSuccessfulSyncAge.ID).
				Msg("failed to parse connection sync age as a duration")
			continue
		}

		connectionsLastSuccessfulSyncHistogramVec.
			WithLabelValues(
				connectionLastSuccessfulSyncAge.DestinationConnector,
				connectionLastSuccessfulSyncAge.SourceConnector,
			).
			Observe(age.Hours())
	}

	connectionsLastSuccessfulSyncHistogramVec.Collect(ch)
}

// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "airbyte"
)

type collector struct{}

// NewCollector initializes and returns a Prometheus collector for Airbyte metrics.
func NewCollector() *collector {
	return &collector{}
}

// Describe publishes the description of each Airbyte metric to a metrics
// channel.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {}

// Collect gathers metrics from Airbyte.
func (c *collector) Collect(ch chan<- prometheus.Metric) {}

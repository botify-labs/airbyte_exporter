// Copyright 2023 VirtualTam.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/virtualtam/airbyte_exporter/internal/airbyte"
)

const (
	webroot = `<html>
<head><title>Airbyte Exporter</title></head>
<body>
  <h1>Airbyte Exporter</h1>
  <p><a href="/metrics">Metrics</a></p>
</body>
</html>`
)

func accessLogger(r *http.Request, status, size int, dur time.Duration) {
	hlog.FromRequest(r).Info().
		Dur("duration_ms", dur).
		Str("host", r.Host).
		Str("path", r.URL.Path).
		Int("size", size).
		Int("status", status).
		Msg("handle request")
}

func newServer(airbyteService *airbyte.Service, listenAddr string) *http.Server {
	collector := NewCollector(airbyteService)
	prometheus.MustRegister(collector)

	router := http.NewServeMux()

	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(webroot))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Setup structured logging middleware
	chain := alice.New(hlog.NewHandler(log.Logger), hlog.AccessHandler(accessLogger))

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      chain.Then(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return server
}

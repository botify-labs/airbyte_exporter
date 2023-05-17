# Airbyte Prometheus Exporter

<img src="https://github.com/virtualtam/airbyte_exporter/actions/workflows/ci.yaml/badge.svg?branch=main" alt="Continuous integration workflow status">
<img src="https://github.com/virtualtam/airbyte_exporter/actions/workflows/docker.yaml/badge.svg?branch=main" alt="Docker image workflow status">

## Metrics exposed
### Counters
- `airbyte_connections_total{destination_connector, source_connector, status}`
- `airbyte_jobs_completed_total{destination_connector, source_connector, type, status}`

### Gauges
- `airbyte_jobs_pending{destination_connector, source_connector, type}`
- `airbyte_jobs_running{destination_connector, source_connector, type}`

## Change Log
See [CHANGELOG](./CHANGELOG.md)

## License
`airbyte_exporter` is licensed under the MIT License.

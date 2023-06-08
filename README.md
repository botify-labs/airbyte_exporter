# Airbyte Prometheus Exporter

<img src="https://github.com/botify-labs/airbyte_exporter/actions/workflows/ci.yaml/badge.svg?branch=main" alt="Continuous integration workflow status">
<img src="https://github.com/botify-labs/airbyte_exporter/actions/workflows/docker.yaml/badge.svg?branch=main" alt="Docker image workflow status">

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/botify-helm-charts)](https://artifacthub.io/packages/search?repo=botify-helm-charts)

## Metrics exposed
### Counters
- `airbyte_connections_total{destination_connector, source_connector, status}`
- `airbyte_jobs_completed_total{destination_connector, source_connector, type, status}`

### Gauges
- `airbyte_jobs_pending{destination_connector, source_connector, type}`
- `airbyte_jobs_running{destination_connector, source_connector, type}`

## Configuration
`airbyte_exporter` can be configured via:

- environment variables, e.g. `AIRBYTE_EXPORTER_DB_PASSWORD=p455w0rd`
- a configuration file
- POSIX flags, e.g. `--db-password p455w0rd`

Available flags can be listed using the program's help:

```shell
$ ./airbyte_exporter --help

Airbyte Exporter

Usage:
  airbyte_exporter [flags]

Flags:
      --db-addr string       Database address (host:port) (default "localhost:5432")
      --db-name string       Database name (default "airbyte")
      --db-password string   Database password (default "airbyte_exporter")
      --db-sslmode string    Database sslmode (default "disable")
      --db-user string       Database user (default "airbyte_exporter")
  -h, --help                 help for airbyte_exporter
      --listen-addr string   Listen to this address (host:port) (default "0.0.0.0:8080")
      --log-level string     Log level (trace, debug, info, warn, error, fatal, panic) (default "info")
```

### Example configuration file

The exporter will look for configuration files located under:
    - `/etc/airbyte_exporter.yaml`
    - `~/.config/airbyte_exporter.yaml`

```yaml
# Global exporter options
listen-addr: 0.0.0.0:8080
log-level: info

# Airbyte database options
db-addr: "postgresql:5432"
db-name: airbyte
db-password: "ch4ng3m3!"
db-sslmode: require
```

### PostgreSQL user
The exporter needs to be able to connect to the Airbyte database, and have read-only access
to Airbyte database tables.

The following commands are provided as an example; see PostgreSQL's documentation for
further information:

- [`CREATE ROLE`](https://www.postgresql.org/docs/current/sql-createrole.html)
- [`GRANT`](https://www.postgresql.org/docs/current/sql-grant.html)
- [Predefined roles](https://www.postgresql.org/docs/current/predefined-roles.html#PREDEFINED-ROLES-TABLE)

Connect to the `airbyte` database:

```shell
# psql -h <host> -p <port> -U <admin_user> <database>
$ psql -h 127.0.0.1 -p 5432 -U postgres airbyte
```

Create the user:

```sql
CREATE ROLE airbyte_exporter WITH LOGIN ENCRYPTED PASSWORD 'SomeStrongPassword';
```

For PostgreSQL **version 14 and above**, grant read-only privileges with:

```sql
GRANT pg_read_all_data TO airbyte_exporter;
```

For PostgreSQL **version 13 and below**, grant read-only privileges on current and newly created tables with:

```sql
-- Current tables
GRANT CONNECT ON DATABASE YourDatabaseName TO airbyte_exporter;
GRANT USAGE ON SCHEMA public TO airbyte_exporter;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO airbyte_exporter;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO airbyte_exporter;
REVOKE CREATE ON SCHEMA public FROM PUBLIC;

-- Newly created tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO airbyte_exporter;
```

## Running

### Running with Docker
See [`airbyte_exporter` container packages](https://github.com/botify-labs/airbyte_exporter/pkgs/container/airbyte_exporter)
for a list of available Docker image tags.

Pull the Docker image:

```shell
$ docker pull ghcr.io/botify-labs/airbyte_exporter:latest
```

Run the exporter:

```shell
$ docker run \
    --name airbyte-exporter \
    --rm \
    -e AIRBYTE_EXPORTER_DB_HOST=postgresql \
    -e AIRBYTE_EXPORTER_DB_PASSWORD=ch4ng3m3 \
    -p 8080:8080
    ghcr.io/botify-labs/airbyte_exporter:latest
```

### Deploying to Kubernetes with Helm

See instructions on Artifact Hub for [botify-helm-charts/prometheus-airbyte-exporter](https://artifacthub.io/packages/helm/botify-helm-charts/prometheus-airbyte-exporter).

## Change Log
See [CHANGELOG](./CHANGELOG.md)

## Contributing
See [CONTRIBUTING](./CONTRIBUTING.md)

## License
`airbyte_exporter` is licensed under the MIT License.

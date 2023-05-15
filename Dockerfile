# Step 1: Build Go binaries
FROM golang:1.20-bullseye as builder

ARG CGO_ENABLED=1

WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

ADD . .
RUN --mount=type=cache,target=/root/.cache/go-build make build

# Step 2: Build the actual image
FROM debian:bullseye-slim

RUN groupadd \
        --gid 1000 \
        airbyte_exporter \
    && useradd \
        --create-home \
        --home-dir /var/lib/airbyte_exporter \
        --shell /bin/bash \
        --uid 1000 \
        --gid airbyte_exporter \
        airbyte_exporter

COPY --from=builder /app/build/airbyte_exporter /usr/local/bin/airbyte_exporter

ENV \
    AIRBYTE_EXPORTER_DB_ADDR="postgres:5432" \
    AIRBYTE_EXPORTER_DB_SSLMODE="disable" \
    AIRBYTE_EXPORTER_DB_NAME="airbyte" \
    AIRBYTE_EXPORTER_DB_USER="airbyte_exporter" \
    AIRBYTE_EXPORTER_DB_PASSWORD="airbyte_exporter" \
    AIRBYTE_EXPORTER_LISTEN_ADDR="0.0.0.0:8080" \
    AIRBYTE_EXPORTER_LOG_LEVEL="info"

EXPOSE 8080

USER airbyte_exporter
WORKDIR /var/lib/airbyte_exporter

CMD ["/usr/local/bin/airbyte_exporter"]

# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [v2.0.1](https://github.com/botify-labs/airbyte_exporter/releases/tag/v2.0.1) - 2023-06-12

### Changed

- Breaking: update module imports to `/v2`

## [v2.0.0](https://github.com/botify-labs/airbyte_exporter/releases/tag/v2.0.0) - 2023-06-08

### Changed

- Breaking: transfer repository to `botify-labs` and update Go module imports
- Update documentation

## [v1.1.0](https://github.com/botify-labs/airbyte_exporter/releases/tag/v1.1.0) - 2023-05-17

### Fixed

- Update SQL queries to perform inner joins between Airbyte tables

## [v1.0.0](https://github.com/botify-labs/airbyte_exporter/releases/tag/v1.0.0) - 2023-05-17

Initial release.

### Added

- Setup exporter repository
- Setup continuous integration with Github Actions
- Publish Docker images for the `main` branch and tags to `ghcr.io`
- Gather Airbyte metrics by querying its PostgreSQL database

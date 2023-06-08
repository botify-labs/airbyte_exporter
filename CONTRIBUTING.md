# Contributing
## Building

Get the sources:

```shell
$ git clone https://github.com/botify-labs/airbyte_exporter.git
$ cd airbyte_exporter
```

Run linters:

```shell
$ make lint
```

Build the parser and exporter:

```shell
$ make build
```

Build platform-specific binaries with [Promu](https://github.com/prometheus/promu):

```shell
$ promu crossbuild
```

Build and archive platform-specific binaries:

```shell
$ promu crossbuild
$ promu crossbuild tarballs
```

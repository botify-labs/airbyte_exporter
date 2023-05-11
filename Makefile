# A common Makefile that includes rules to be reused in different Prometheus projects.
# https://github.com/prometheus/prometheus/blob/master/Makefile.common
include Makefile.common

all: lint cover build
.PHONY: all

clean:
	rm -rf .build .tarballs airbyte_exporter
.PHONY: clean

lint:
	golangci-lint run ./...
.PHONY: lint

test:
	go test -race ./...
.PHONY: test

cover:
	go test -cover -race ./...
.PHONY: cover

release: clean
	promu crossbuild
	promu crossbuild tarballs
	cd .tarballs; sha256sum * > sha256sums
.PHONY: release

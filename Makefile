APP_IMAGE_NAME      = currency-loader
REF_NAME           ?= $(shell git rev-parse --abbrev-ref HEAD)
IMAGE_VERSION      ?= ${REF_NAME}-$(shell git rev-parse HEAD)

export IMAGE_VERSION

SRCDIR1 := /go/src/gitlab.2gis.ru/traffic/tugc-aggregator
DOCKERFLAGS := --rm --user $$(id -u):$$(id -g) -v $(CURDIR):$(SRCDIR1):rw -w $(SRCDIR1)

PROJECT_PKGS := $$(go list ./...)

.PHONY: clean-api
clean-api:
	rm -rf bin/currency-loader/*

.PHONY: build-api
build-api: clean-api
	CGO_ENABLED=0 GOOS=linux go build -o bin/currency-loader/currency-loader ./src/cmd

.PHONY: build-app
build-app:
	docker-compose run --rm currency-loader-build

.PHONY: build-app-image
build-app-image:
	docker-compose build currency-loader

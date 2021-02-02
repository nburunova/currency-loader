APP_IMAGE_NAME      = currency-loader
REF_NAME           ?= $(shell git rev-parse --abbrev-ref HEAD)
IMAGE_VERSION      ?= ${REF_NAME}-$(shell git rev-parse HEAD)

export IMAGE_VERSION

.PHONY: clean-api
clean-api:
	rm -rf bin/currency-loader/*

.PHONY: build-api
build-api: clean-api
	go build -a -installsuffix cgo -ldflags "-w -s" -o bin/currency-loader/currency-loader ./src/cmd

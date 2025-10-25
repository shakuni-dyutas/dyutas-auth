include .env

.DEFAULT_GOAL := build

og: # api code generation
	oapi-codegen -config api-spec/oapi-codegen.yaml api-spec/openapi.yaml
.PHONY: og

svt: # sql code vet
	AUTH_RDB_PASSWORD=$(AUTH_RDB_PASSWORD) AUTH_RDB_PORT=$(AUTH_RDB_PORT) AUTH_RDB_DB=$(AUTH_RDB_DB) sqlc vet
.PHONY: svt

sg: svt # sql code generation
	AUTH_RDB_PASSWORD=$(AUTH_RDB_PASSWORD) AUTH_RDB_PORT=$(AUTH_RDB_PORT) AUTH_RDB_DB=$(AUTH_RDB_DB) sqlc generate
.PHONY: sg

go: # run the app with air for hot-reload. Just use vscode launch task instead of this.
	air
.PHONY: go

fmt:
	go fmt ./...
.PHONY: fmt

lint: fmt
	golint ./...
.PHONY: lint

vet: fmt
	go vet ./...
	# shadow ./... # this tool detects shadowing variables
.PHONY: vet

build: vet
	go build
.PHONY: build
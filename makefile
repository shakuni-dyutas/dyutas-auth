include .env
include $(DYUTAS_ENV_PATH)

.DEFAULT_GOAL := build

dg: # generate package dependency graph
	./dev_env/godep_graph.sh
.PHONY: dg

og: # api code generation
	oapi-codegen -config api-spec/oapi-codegen.yaml api-spec/openapi.yaml
.PHONY: og

svt: # sql code vet
	source ./dev_env/load_env_vars.sh && sqlc vet
.PHONY: svt

sg: svt # sql code generation
	source ./dev_env/load_env_vars.sh && sqlc generate
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
	go build -o ./tmp/dyutas-auth ./cmd/dyutas-auth
.PHONY: build
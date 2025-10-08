.DEFAULT_GOAL := build

go: # run the app with air for hot-reload
	air
.PHONY:go

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
	# shadow ./... # this tool detects shadowing variables
.PHONY:vet

build: vet
	go build
.PHONY:build
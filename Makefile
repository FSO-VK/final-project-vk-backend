.PHONY: setup dev format lint test

setup:
	go mod download
	go tool lefthook install

dev:
	docker compose -f compose.dev.yml up --build --watch

format:
	golangci-lint fmt

lint:
	golangci-lint run

test:
	# there is a bug in GOTOOLCHAIN with go 1.25.x
	GOTOOLCHAIN=go1.25.0+auto gotestsum --format pkgname -- -race -coverprofile=coverage.out ./...
.PHONY: format lint test

dev:
	docker compose -f compose.dev.yml up --watch

format:
	golangci-lint fmt

lint:
	golangci-lint run

test:
	gotestsum --format pkgname -- -race -coverprofile=coverage.out ./...
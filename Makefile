.PHONY: format lint test

format:
	golangci-lint fmt

lint:
	golangci-lint run

test:
	gotestsum --format pkgname -- -race -coverprofile=coverage.out ./...
	rm coverage.out
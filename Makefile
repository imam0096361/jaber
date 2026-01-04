.PHONY: help run build clean test dev prod

help:
	@echo "News Portal - Available commands:"
	@echo "  make run      - Run development server"
	@echo "  make build    - Build binary"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make dev      - Run with hot reload"
	@echo "  make prod     - Build for production"

run:
	go run src/main.go

build:
	go build -o news-portal src/main.go

clean:
	rm -f news-portal
	rm -f news.db
	rm -rf static/uploads/*

test:
	go test ./...

dev:
	@command -v air > /dev/null || go install github.com/cosmtrek/air@latest
	air

prod:
	CGO_ENABLED=1 go build -ldflags="-s -w" -o news-portal src/main.go
	@echo "Build complete: ./news-portal"

install-deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

lint:
	@command -v golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

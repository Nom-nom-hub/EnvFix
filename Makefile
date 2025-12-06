.PHONY: build clean test run help

# Build targets
build:
	go build -o envfix.exe -v

build-linux:
	GOOS=linux GOARCH=amd64 go build -o envfix -v

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o envfix -v

build-all: build build-linux build-macos

# Run targets
run:
	go run main.go

run-scan:
	go run main.go scan

run-repair:
	go run main.go repair --dry-run

run-clean:
	go run main.go clean

run-lock:
	go run main.go lock

run-doctor:
	go run main.go doctor

# Testing
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Code quality
fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

# Dependencies
deps:
	go mod tidy
	go mod download

# Cleaning
clean:
	go clean
	rm -f envfix envfix.exe coverage.out

# Help
help:
	@echo "Available targets:"
	@echo "  make build          - Build envfix for Windows"
	@echo "  make build-linux    - Build envfix for Linux"
	@echo "  make build-macos    - Build envfix for macOS"
	@echo "  make build-all      - Build for all platforms"
	@echo "  make run            - Run envfix"
	@echo "  make run-scan       - Run scan command"
	@echo "  make run-repair     - Run repair command (dry-run)"
	@echo "  make run-clean      - Run clean command"
	@echo "  make run-lock       - Run lock command"
	@echo "  make run-doctor     - Run doctor command"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make vet            - Run vet"
	@echo "  make deps           - Download dependencies"
	@echo "  make clean          - Clean build artifacts"

.PHONY: build run test clean deps lint

# Build the application
build:
	@go build -o ./bin/api ./cmd/api

dev:
	@gin --appPort 8080 --all run ./main.go

# Run the application
run:
	@go run ./main.go

# Run tests
test:
	@go test -v ./...

# Clean build artifacts
clean:
	rm -rf ./bin

# Download and install dependencies
deps:
	@go mod download
	@go mod tidy

# Lint the code
lint:
	@golangci-lint run
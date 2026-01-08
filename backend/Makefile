.PHONY: run migrate test build

# Run the API server
run:
	go run cmd/main.go --api --port=8080

# Run database migrations
migrate:
	go run cmd/main.go --migrate

# Run tests
test:
	go test ./...

# Build the application
build:
	go build -o bin/oms cmd/main.go

# Install dependencies
deps:
	go mod download
	go mod tidy


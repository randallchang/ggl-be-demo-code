.PHONY: init build run test clean docker-build docker-run mock fmt lint tidy

# Initialize project
init:
	go mod download
	go generate ./...

# Build project
build:
	go build -o app ./cmd/main

# Run project
run:
	go run ./cmd/main

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f app
	go clean

# Build Docker image
docker-build:
	docker build -f build/Dockerfile -t task-service .

# Run Docker container
docker-run:
	docker run -p 8080:8080 task-service

# Generate mock files
mock:
	go generate ./...

# Format code
fmt:
	go fmt ./...

# Run code checks
lint:
	go vet ./...

# Update dependencies
tidy:
	go mod tidy
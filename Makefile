.PHONY: build run test clean deps migrate

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run with hot reload (requires air)
dev:
	air

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Generate documentation
docs:
	swag init -g cmd/api/main.go

# Database migration (manual - run SQL script)
migrate:
	@echo "Please run the SQL script in migrations/001_create_users_table.sql manually in your Neon database"

# Help
help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install dependencies"
	@echo "  dev        - Run with hot reload (requires air)"
	@echo "  build-prod - Build for production"
	@echo "  lint       - Run linter"
	@echo "  fmt        - Format code"
	@echo "  docs       - Generate documentation"
	@echo "  migrate    - Show migration instructions" 
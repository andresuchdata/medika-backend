# Medika Backend Makefile

.PHONY: help build run test docker-up docker-down migrate clean deps lint

# Default target
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application in development mode"
	@echo "  test        - Run tests"
	@echo "  docker-up   - Start all services with Docker Compose"
	@echo "  docker-down - Stop all Docker services"
	@echo "  migrate     - Run database migrations"
	@echo "  seed        - Seed database with all test data"
	@echo "  setup-db    - Run migrations and seed data"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Download dependencies"
	@echo "  lint        - Run linter"

# Build the application
build:
	@echo "Building medika-backend..."
	@go build -o bin/api cmd/api/main.go
	@go build -o bin/seeder cmd/seeder/main.go

# Run the application
run:
	@echo "Starting medika-backend..."
	@go run cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Start Docker services
docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down

# Build and start all services
docker-build:
	@echo "Building and starting all services..."
	@docker-compose up --build -d

# Run database migrations
migrate:
	@echo "Running database migrations..."
	@PGPASSWORD=medika_pass psql -h localhost -U medika_user -d medika_db -f migrations/001_initial_schema.sql

# Seed database with test data
seed:
	@echo "Seeding database..."
	@go run cmd/seeder/main.go -all

# Seed specific data
seed-organizations:
	@echo "Seeding organizations..."
	@go run cmd/seeder/main.go -organizations

seed-users:
	@echo "Seeding users..."
	@go run cmd/seeder/main.go -users

seed-rooms:
	@echo "Seeding rooms..."
	@go run cmd/seeder/main.go -rooms

# Full setup: migrate + seed
setup-db: migrate seed
	@echo "Database setup completed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/cosmtrek/air@latest

# Generate swagger docs
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go

# Run development server with hot reload
dev:
	@echo "Starting development server with hot reload..."
	@air

# Check security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	@govulncheck ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run all checks
check: fmt lint test security
	@echo "All checks completed!"

# Docker logs
logs:
	@docker-compose logs -f

# Database shell
db-shell:
	@PGPASSWORD=medika_pass psql -h localhost -U medika_user -d medika_db

# Redis shell
redis-shell:
	@docker-compose exec redis redis-cli

# Environment setup
env-setup:
	@echo "Setting up environment..."
	@cp .env.example .env
	@echo "Please edit .env file with your configuration"

# Production build
build-prod:
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go

# Local development with hot reload and dependencies
dev-full: docker-up
	@sleep 5  # Wait for services to start
	@make migrate
	@make dev

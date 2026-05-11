.PHONY: help install build run test clean dev serve migrate seed logs

help:
	@echo "Govel Framework - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make install      Install dependencies"
	@echo "  make dev          Start development server with hot reload"
	@echo "  make serve        Start development server"
	@echo ""
	@echo "Building:"
	@echo "  make build        Build binary"
	@echo "  make build-release Build optimized release binary"
	@echo ""
	@echo "Testing:"
	@echo "  make test         Run tests"
	@echo "  make test-coverage Run tests with coverage report"
	@echo ""
	@echo "Database:"
	@echo "  make migrate      Run database migrations"
	@echo "  make migrate-down Rollback migrations"
	@echo "  make seed         Seed database"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt          Format code"
	@echo "  make lint         Run linter"
	@echo "  make vet          Run go vet"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean        Remove build artifacts"
	@echo "  make logs         View application logs"

install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

dev:
	@echo "Starting development server..."
	go run main.go

serve:
	@echo "Starting server..."
	go run main.go

build:
	@echo "Building application..."
	go build -o app main.go

build-release:
	@echo "Building release binary..."
	go build -ldflags "-s -w" -o app main.go

run: build
	@echo "Running application..."
	./app

test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run ./...

vet:
	@echo "Running go vet..."
	go vet ./...

clean:
	@echo "Cleaning up..."
	rm -f app
	rm -f coverage.out coverage.html
	go clean

migrate:
	@echo "Running migrations..."
	go run artisan.go migrate

migrate-down:
	@echo "Rollback migrations..."
	go run artisan.go migrate:rollback

seed:
	@echo "Seeding database..."
	go run artisan.go db:seed

logs:
	@echo "Showing application logs..."
	tail -f storage/logs/*.log

# Make commands
make-controller:
	@echo "Creating controller..."
	go run artisan.go make:controller

make-model:
	@echo "Creating model..."
	go run artisan.go make:model

make-migration:
	@echo "Creating migration..."
	go run artisan.go make:migration

make-middleware:
	@echo "Creating middleware..."
	go run artisan.go make:middleware

make-request:
	@echo "Creating form request..."
	go run artisan.go make:request

make-event:
	@echo "Creating event..."
	go run artisan.go make:event

make-listener:
	@echo "Creating listener..."
	go run artisan.go make:listener

make-job:
	@echo "Creating job..."
	go run artisan.go make:job

# Docker support (optional)
docker-build:
	docker build -t govel .

docker-run:
	docker run -p 8000:8000 govel

# Pre-commit checks
check: fmt vet lint test
	@echo "All checks passed!"

# Full setup for new developers
setup: install migrate seed
	@echo "Setup complete!"

.DEFAULT_GOAL := help

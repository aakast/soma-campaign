.PHONY: help build run dev clean docker-build docker-up docker-down test

# Default target
help:
	@echo "Soma Mayel Campaign Website - Available Commands:"
	@echo ""
	@echo "  make build         - Build the Go application"
	@echo "  make run           - Run the application locally"
	@echo "  make dev           - Run in development mode with hot reload"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start with Docker Compose"
	@echo "  make docker-down   - Stop Docker Compose"
	@echo "  make test          - Run tests"
	@echo "  make setup         - Initial setup (install dependencies)"

# Setup environment
setup:
	@echo "Setting up environment..."
	@cp -n .env.example .env || true
	@go mod download
	@echo "Setup complete! Edit .env file with your configuration."

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/soma-campaign main.go
	@echo "Build complete! Binary available at bin/soma-campaign"

# Run the application
run: build
	@echo "Starting application..."
	@./bin/soma-campaign

# Development mode with auto-reload (requires air)
dev:
	@echo "Starting in development mode..."
	@if ! command -v air &> /dev/null; then \
		echo "Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@air

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f soma-campaign
	@echo "Clean complete!"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker-compose build
	@echo "Docker build complete!"

docker-up:
	@echo "Starting with Docker Compose..."
	@docker-compose up -d
	@echo "Application running at http://localhost:3000"

docker-down:
	@echo "Stopping Docker Compose..."
	@docker-compose down
	@echo "Docker compose stopped!"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...
	@echo "Tests complete!"

# Create necessary directories
init-dirs:
	@mkdir -p content/posts content/pages static/images static/videos
	@echo "Directories created!"

# Production deployment
deploy: docker-build
	@echo "Deploying to production..."
	@docker-compose -f docker-compose.yml up -d
	@echo "Deployment complete!"

# View logs
logs:
	@docker-compose logs -f web

# Database backup (if needed in future)
backup:
	@echo "Creating backup..."
	@tar -czf backup-$(shell date +%Y%m%d-%H%M%S).tar.gz content/ static/images/ static/videos/
	@echo "Backup complete!"
# B3 Challenge Makefile
# Provides convenient commands for development, testing, and deployment

.PHONY: help build test clean run docker-build docker-run setup-db load-data benchmark

# Default target
help: ## Show this help message
	@echo "B3 Challenge - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make setup     # Complete project setup"
	@echo "  make run       # Start the API server"
	@echo "  make test      # Run all tests"
	@echo "  make benchmark # Run performance tests"

# Variables
APP_NAME=b3-challenge
API_PORT=8080
DB_NAME=b3_challenge
DB_USER=b3user
DB_PASSWORD=b3password

# Build targets
build: ## Build the application binaries
	@echo "Building application binaries..."
	@go build -o bin/api cmd/api/main.go
	@go build -o bin/data-loader cmd/data-loader/main.go
	@echo "✓ Build completed successfully"

build-linux: ## Build for Linux (useful for Docker)
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o bin/api-linux cmd/api/main.go
	@GOOS=linux GOARCH=amd64 go build -o bin/data-loader-linux cmd/data-loader/main.go
	@echo "✓ Linux build completed"

# Development targets
run: ## Start the API server
	@echo "Starting B3 Challenge API server..."
	@go run cmd/api/main.go

run-data-loader: ## Run the data loader
	@echo "Running data loader..."
	@go run cmd/data-loader/main.go

# Testing targets
test: ## Run all tests
	@echo "Running tests..."
	@go test ./... -v

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test ./... -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test ./internal/service ./internal/api -v

benchmark: ## Run performance benchmarks
	@echo "Running performance benchmarks..."
	@go test ./internal/service -bench=. -benchmem
	@go test ./internal/api -bench=. -benchmem

benchmark-api: ## Run API performance test (requires running server)
	@echo "Running API performance test..."
	@./scripts/simple_performance_test.sh

# Database targets
setup-db: ## Setup PostgreSQL database
	@echo "Setting up PostgreSQL database..."
	@sudo systemctl start postgresql || true
	@sudo -u postgres psql -c "CREATE DATABASE $(DB_NAME);" 2>/dev/null || echo "Database already exists"
	@sudo -u postgres psql -c "CREATE USER $(DB_USER) WITH PASSWORD '$(DB_PASSWORD)';" 2>/dev/null || echo "User already exists"
	@sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $(DB_NAME) TO $(DB_USER);" 2>/dev/null || true
	@echo "✓ Database setup completed"

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .
	@echo "✓ Docker image built: $(APP_NAME):latest"

docker-run: ## Run application in Docker
	@echo "Running application in Docker..."
	@docker run -p $(API_PORT):$(API_PORT) --env-file .env $(APP_NAME):latest

docker-compose-up: ## Start application with Docker Compose
	@echo "Starting application with Docker Compose..."
	@docker compose up -d postgres
	@docker compose up -d
	@echo "✓ Application started at http://localhost:$(API_PORT)"

docker-compose-down: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	@docker compose down
	@echo "✓ Services stopped"

# Setup and installation targets
setup: ## Complete project setup
	@echo "Setting up B3 Challenge project..."
	@docker compose up -d postgres
	@echo "waiting for PostgreSQL to start..."
	@sleep 5
	@echo "1. Installing Go dependencies..."
	@go mod tidy
	@echo "2. Setting up database..."
	@make setup-db
	@echo "3. Creating environment file..."
	@cp .env.example .env 2>/dev/null || echo ".env already exists"
	@echo "4. Building application..."
	@make build
	@echo "5. Loading sample data..."
	@make run-data-loader
	@echo ""
	@echo "✓ Setup completed successfully!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Review and update .env file if needed"
	@echo "  2. Run 'make run' to start the API server"
	@echo "  3. Test the API with 'curl http://localhost:$(API_PORT)/health'"

# Maintenance targets
clean: ## Clean build artifacts and temporary files
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean -cache
	@echo "✓ Cleanup completed"

fmt: ## Format Go code
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "✓ Code formatted"

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@golangci-lint run ./... || echo "Install golangci-lint for linting: https://golangci-lint.run/usage/install/"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet completed"

status: ## Check application status
	@echo "Checking application status..."
	@curl -s http://localhost:$(API_PORT)/health | jq . || echo "API not responding"

deps: ## Show dependency tree
	@echo "Dependency tree:"
	@go mod graph

update-deps: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✓ Dependencies updated"

ci: ## Run CI pipeline locally
	@echo "Running CI pipeline..."
	@make fmt
	@make vet
	@make test
	@make benchmark
	@make build
	@echo "✓ CI pipeline completed successfully"

# Production targets
prod-build: ## Build for production
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/api cmd/api/main.go
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/data-loader cmd/data-loader/main.go
	@echo "✓ Production build completed"

deploy: ## Deploy to production (customize as needed)
	@echo "Deploying to production..."
	@echo "This target should be customized for your deployment environment"
	@echo "Consider using Docker, Kubernetes, or your preferred deployment method"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@swag init --generalInfo cmd/api/main.go --output docs
	@echo "✓ Swagger documentation generated in docs/"
	@echo "You can view it at http://localhost:$(API_PORT)/swagger/index.html"

# Help target (default)
.DEFAULT_GOAL := help


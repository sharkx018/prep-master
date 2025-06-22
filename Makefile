.PHONY: help dev prod build up down logs clean

# Default target
help:
	@echo "Available commands:"
	@echo "  make dev       - Start development environment with hot reload"
	@echo "  make prod      - Start production environment"
	@echo "  make build     - Build all Docker images"
	@echo "  make up        - Start all services"
	@echo "  make down      - Stop all services"
	@echo "  make logs      - View logs from all services"
	@echo "  make clean     - Clean up volumes and containers"
	@echo "  make db-shell  - Access PostgreSQL shell"
	@echo "  make backend-shell - Access backend container shell"
	@echo "  make frontend-shell - Access frontend container shell"

# Development environment with hot reload
dev:
	docker-compose -f docker-compose.dev.yml up

# Production environment
prod:
	docker-compose up

# Build all images
build:
	docker-compose build

# Start services in background
up:
	docker-compose up -d

# Stop services
down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# Clean everything
clean:
	docker-compose down -v
	docker system prune -f

# Database shell
db-shell:
	docker exec -it interview-prep-db psql -U interview_user -d interview_prep

# Backend shell
backend-shell:
	docker exec -it interview-prep-backend /bin/sh

# Frontend shell
frontend-shell:
	docker exec -it interview-prep-frontend /bin/sh

# Run backend tests
test-backend:
	cd backend && go test ./...

# Run frontend tests
test-frontend:
	cd frontend && npm test

# Format Go code
fmt:
	cd backend && go fmt ./...

# Lint Go code
lint:
	cd backend && golangci-lint run 
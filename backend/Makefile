# Interview Prep App Makefile

# Variables
DOCKER_COMPOSE = docker-compose
GO = go
APP_NAME = interview-prep-app

# Colors for output
GREEN = \033[0;32m
YELLOW = \033[0;33m
RED = \033[0;31m
NC = \033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${GREEN}%-20s${NC} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Docker commands
.PHONY: docker-up
docker-up: ## Start all services with Docker Compose
	@echo "${YELLOW}Starting Docker services...${NC}"
	$(DOCKER_COMPOSE) up -d
	@echo "${GREEN}Services started! App running at http://localhost:8080${NC}"

.PHONY: docker-down
docker-down: ## Stop all services
	@echo "${YELLOW}Stopping Docker services...${NC}"
	$(DOCKER_COMPOSE) down
	@echo "${GREEN}Services stopped!${NC}"

.PHONY: docker-restart
docker-restart: docker-down docker-up ## Restart all services

.PHONY: docker-logs
docker-logs: ## Show logs from all services
	$(DOCKER_COMPOSE) logs -f

.PHONY: docker-logs-app
docker-logs-app: ## Show logs from app service only
	$(DOCKER_COMPOSE) logs -f app

.PHONY: docker-logs-db
docker-logs-db: ## Show logs from database service only
	$(DOCKER_COMPOSE) logs -f postgres

.PHONY: docker-build
docker-build: ## Rebuild Docker images
	@echo "${YELLOW}Building Docker images...${NC}"
	$(DOCKER_COMPOSE) build
	@echo "${GREEN}Build complete!${NC}"

.PHONY: docker-clean
docker-clean: ## Remove all containers and volumes (WARNING: This will delete all data!)
	@echo "${RED}WARNING: This will delete all containers and volumes!${NC}"
	@read -p "Are you sure? (y/N) " confirm && [ "$$confirm" = "y" ] || exit 1
	$(DOCKER_COMPOSE) down -v
	@echo "${GREEN}Cleanup complete!${NC}"

# Database commands
.PHONY: db-only
db-only: ## Start only the PostgreSQL database
	@echo "${YELLOW}Starting PostgreSQL database...${NC}"
	$(DOCKER_COMPOSE) up -d postgres
	@echo "${GREEN}Database started! Connection: postgresql://interview_user:interview_pass@localhost:5432/interview_prep${NC}"

.PHONY: db-shell
db-shell: ## Open PostgreSQL shell
	@echo "${YELLOW}Connecting to PostgreSQL...${NC}"
	$(DOCKER_COMPOSE) exec postgres psql -U interview_user -d interview_prep

.PHONY: db-backup
db-backup: ## Backup the database
	@echo "${YELLOW}Creating database backup...${NC}"
	@mkdir -p backups
	$(DOCKER_COMPOSE) exec postgres pg_dump -U interview_user interview_prep > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "${GREEN}Backup created in backups/ directory${NC}"

.PHONY: db-restore
db-restore: ## Restore database from backup (Usage: make db-restore FILE=backups/backup_file.sql)
	@if [ -z "$(FILE)" ]; then echo "${RED}Please specify backup file: make db-restore FILE=backups/backup_file.sql${NC}"; exit 1; fi
	@echo "${YELLOW}Restoring database from $(FILE)...${NC}"
	$(DOCKER_COMPOSE) exec -T postgres psql -U interview_user interview_prep < $(FILE)
	@echo "${GREEN}Database restored!${NC}"

# Development commands
.PHONY: dev
dev: ## Run the app locally (requires local PostgreSQL or use 'make db-only' first)
	@echo "${YELLOW}Starting app in development mode...${NC}"
	@cp env.docker .env 2>/dev/null || true
	$(GO) run cmd/server/main.go

.PHONY: build
build: ## Build the Go application
	@echo "${YELLOW}Building application...${NC}"
	$(GO) build -o $(APP_NAME) cmd/server/main.go
	@echo "${GREEN}Build complete! Binary: $(APP_NAME)${NC}"

.PHONY: test
test: ## Run tests
	@echo "${YELLOW}Running tests...${NC}"
	$(GO) test ./...

.PHONY: lint
lint: ## Run linter
	@echo "${YELLOW}Running linter...${NC}"
	golangci-lint run

.PHONY: mod-tidy
mod-tidy: ## Tidy Go modules
	@echo "${YELLOW}Tidying Go modules...${NC}"
	$(GO) mod tidy

# Quick start commands
.PHONY: setup
setup: ## First time setup - start database and copy env file
	@echo "${YELLOW}Setting up Interview Prep App...${NC}"
	@cp env.docker .env 2>/dev/null || echo "${YELLOW}.env file already exists${NC}"
	@$(MAKE) docker-up
	@echo ""
	@echo "${GREEN}Setup complete!${NC}"
	@echo ""
	@echo "App is running at: ${GREEN}http://localhost:8080${NC}"
	@echo "Database is at: ${GREEN}postgresql://interview_user:interview_pass@localhost:5432/interview_prep${NC}"
	@echo ""
	@echo "Next steps:"
	@echo "  - Check logs: ${YELLOW}make docker-logs${NC}"
	@echo "  - Stop services: ${YELLOW}make docker-down${NC}"
	@echo "  - Connect to DB: ${YELLOW}make db-shell${NC}"

.PHONY: quick-start
quick-start: setup ## Alias for setup

# Status command
.PHONY: status
status: ## Show status of Docker services
	@echo "${YELLOW}Service Status:${NC}"
	@$(DOCKER_COMPOSE) ps 
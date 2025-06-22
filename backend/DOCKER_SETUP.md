# Docker Setup Guide for Interview Prep App

This guide will help you set up and run the Interview Prep App using Docker.

## Prerequisites

- Docker installed on your machine ([Download Docker](https://www.docker.com/get-started))
- Docker Compose (usually comes with Docker Desktop)
- Make (optional, for easier commands)

## Quick Start

### Option 1: Using Make (Recommended)

```bash
# First time setup - starts everything
make setup

# View all available commands
make help
```

### Option 2: Using Docker Compose Directly

```bash
# Start all services
docker-compose up -d

# Check if services are running
docker-compose ps

# View logs
docker-compose logs -f
```

## Common Operations

### 1. Starting Services

```bash
# Start both database and app
make docker-up
# OR
docker-compose up -d

# Start only the database
make db-only
# OR
docker-compose up -d postgres
```

### 2. Stopping Services

```bash
# Stop all services
make docker-down
# OR
docker-compose down

# Stop and remove all data (WARNING!)
make docker-clean
# OR
docker-compose down -v
```

### 3. Viewing Logs

```bash
# All services
make docker-logs
# OR
docker-compose logs -f

# App only
make docker-logs-app
# OR
docker-compose logs -f app

# Database only
make docker-logs-db
# OR
docker-compose logs -f postgres
```

### 4. Database Operations

```bash
# Connect to PostgreSQL shell
make db-shell
# OR
docker-compose exec postgres psql -U interview_user -d interview_prep

# Create a backup
make db-backup

# Restore from backup
make db-restore FILE=backups/backup_20240101_120000.sql
```

## Development Workflow

### Running App Locally with Docker Database

1. Start only the PostgreSQL database:
   ```bash
   make db-only
   ```

2. Copy the Docker environment file:
   ```bash
   cp env.docker .env
   ```

3. Run the Go application locally:
   ```bash
   make dev
   # OR
   go run cmd/server/main.go
   ```

### Building and Running Everything in Docker

```bash
# Rebuild the app image
make docker-build

# Restart services with new build
make docker-restart
```

## Connection Details

When running with Docker:

- **App URL**: http://localhost:8080
- **Database Connection**: 
  - Host: localhost
  - Port: 5432
  - Database: interview_prep
  - Username: interview_user
  - Password: interview_pass
  - Full URL: `postgresql://interview_user:interview_pass@localhost:5432/interview_prep`

## Troubleshooting

### 1. Port Already in Use

If you get an error about port 5432 or 8080 already being in use:

```bash
# Check what's using the ports
lsof -i :5432
lsof -i :8080

# Stop the conflicting service or change ports in docker-compose.yml
```

### 2. Database Connection Issues

If the app can't connect to the database:

```bash
# Check if database is healthy
docker-compose ps

# View database logs
make docker-logs-db

# Try connecting manually
make db-shell
```

### 3. Permission Issues

If you get permission errors:

```bash
# Make sure Docker daemon is running
docker info

# On Linux, you might need to use sudo or add your user to docker group
sudo usermod -aG docker $USER
# Then logout and login again
```

### 4. Clean Start

If things aren't working, try a clean start:

```bash
# Stop everything and remove volumes
make docker-clean

# Start fresh
make setup
```

## Environment Variables

The app uses these environment variables (defined in `docker-compose.yml`):

- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Application port (default: 8080)
- `NODE_ENV`: Environment mode (development/production)

## Data Persistence

- Database data is stored in a Docker volume named `postgres_data`
- This data persists between container restarts
- To completely remove data, use `docker-compose down -v`

## Production Considerations

For production deployment:

1. Change database credentials in `docker-compose.yml`
2. Use environment-specific `.env` files
3. Set up proper backup strategies
4. Configure SSL/TLS for database connections
5. Use a reverse proxy (nginx) for the app
6. Set up monitoring and logging

## Useful Docker Commands

```bash
# List all containers
docker ps -a

# View container logs
docker logs interview_prep_db
docker logs interview_prep_app

# Execute commands in container
docker exec -it interview_prep_db bash
docker exec -it interview_prep_app sh

# Clean up unused resources
docker system prune -a
```

## Next Steps

1. Start the services: `make setup`
2. Test the API: `curl http://localhost:8080/health`
3. Create your first item:
   ```bash
   curl -X POST http://localhost:8080/api/v1/items \
     -H "Content-Type: application/json" \
     -d '{
       "title": "Two Sum",
       "link": "https://leetcode.com/problems/two-sum/",
       "category": "dsa",
       "subcategory": "arrays"
     }'
   ```

Happy coding! 🚀 
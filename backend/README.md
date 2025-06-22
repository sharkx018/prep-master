# Interview Prep App - Backend

Go backend for the Interview Prep application.

## Setup

1. **Install Go** (version 1.21 or higher)

2. **Set up environment variables**
   ```bash
   cp env.example .env
   ```
   Edit `.env` with your database credentials.

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Run the server**
   ```bash
   go run cmd/server/main.go
   ```

## Development

### Run with hot reload
```bash
make dev
```

### Build
```bash
make build
```

### Run tests
```bash
make test
```

### Docker
```bash
docker build -t interview-prep-backend .
docker run -p 8080:8080 --env-file .env interview-prep-backend
```

## API Documentation

The API runs on port 8080 by default. See the main README for endpoint documentation. 
# Interview Prep App - Deployment Guide

This guide covers various deployment options for the Interview Prep application.

## Table of Contents
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Cloud Deployment](#cloud-deployment)
  - [Railway](#railway-recommended)
  - [Vercel + Railway](#vercel--railway)
  - [Render](#render)
  - [Fly.io](#flyio)
  - [DigitalOcean](#digitalocean)

## Local Development

### Without Docker
```bash
# Backend
cd backend
cp env.example .env
go run cmd/server/main.go

# Frontend (in another terminal)
cd frontend
npm install
npm start
```

### With Docker (Recommended)
```bash
# Start everything with hot reload
make dev

# Or manually
docker-compose -f docker-compose.dev.yml up
```

## Docker Deployment

### Production Build
```bash
# Build and start all services
make prod

# Or manually
docker-compose up --build
```

### Docker on VPS
1. **Install Docker and Docker Compose on your VPS**
   ```bash
   curl -fsSL https://get.docker.com -o get-docker.sh
   sh get-docker.sh
   ```

2. **Clone your repository**
   ```bash
   git clone <your-repo-url>
   cd interview-prep-app
   ```

3. **Set up environment**
   ```bash
   cp env.example .env
   # Edit .env with production values
   nano .env
   ```

4. **Start services**
   ```bash
   docker-compose up -d
   ```

5. **Set up Nginx reverse proxy (optional)**
   ```nginx
   server {
       listen 80;
       server_name your-domain.com;
       
       location / {
           proxy_pass http://localhost:3000;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
       
       location /api {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

## Cloud Deployment

### Railway (Recommended)

Railway provides the easiest deployment with automatic SSL, databases, and CI/CD.

#### Full Stack on Railway

1. **Push to GitHub**
   ```bash
   git add .
   git commit -m "Initial commit"
   git push origin main
   ```

2. **Create Railway Project**
   - Go to [railway.app](https://railway.app)
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository

3. **Add Services**
   - Click "New" → "Database" → "Add PostgreSQL"
   - Click "New" → "GitHub Repo" → Select your repo
   - For backend: Set root directory to `/backend`
   - For frontend: Add another service, set root directory to `/frontend`

4. **Configure Backend**
   - Go to backend service → Variables
   - Add:
     ```
     DATABASE_URL=${{Postgres.DATABASE_URL}}
     PORT=8080
     NODE_ENV=production
     GIN_MODE=release
     ```

5. **Configure Frontend**
   - Go to frontend service → Variables
   - Add:
     ```
     REACT_APP_API_URL=https://${{RAILWAY_STATIC_URL}}/api/v1
     ```
   - Go to Settings → Build Command: `npm run build`
   - Start Command: `npm install -g serve && serve -s build -l 3000`

6. **Deploy**
   - Railway will automatically deploy on push

### Vercel + Railway

Use Vercel for frontend (better performance) and Railway for backend.

#### Backend on Railway
Follow steps 1-4 from Railway section above.

#### Frontend on Vercel

1. **Install Vercel CLI**
   ```bash
   npm i -g vercel
   ```

2. **Deploy Frontend**
   ```bash
   cd frontend
   vercel
   ```

3. **Configure Environment**
   - Go to Vercel Dashboard → Settings → Environment Variables
   - Add: `REACT_APP_API_URL` = `https://your-backend.railway.app/api/v1`

4. **Set up Auto-Deploy**
   - Connect GitHub repository
   - Set root directory to `frontend`

### Render

1. **Create Render Account**
   - Sign up at [render.com](https://render.com)

2. **Deploy Database**
   - New → PostgreSQL
   - Choose free tier

3. **Deploy Backend**
   - New → Web Service
   - Connect GitHub repo
   - Root Directory: `backend`
   - Build Command: `go build -o main cmd/server/main.go`
   - Start Command: `./main`
   - Add environment variables

4. **Deploy Frontend**
   - New → Static Site
   - Root Directory: `frontend`
   - Build Command: `npm run build`
   - Publish Directory: `build`

### Fly.io

1. **Install Fly CLI**
   ```bash
   curl -L https://fly.io/install.sh | sh
   ```

2. **Deploy Backend**
   ```bash
   cd backend
   fly launch
   fly postgres create
   fly postgres attach
   fly deploy
   ```

3. **Deploy Frontend**
   ```bash
   cd ../frontend
   fly launch
   fly deploy
   ```

### DigitalOcean

#### Using App Platform

1. **Create App**
   - Go to DigitalOcean → Apps → Create App
   - Connect GitHub repository

2. **Configure Components**
   - Add Database: PostgreSQL
   - Add Service: Backend (`/backend`)
   - Add Static Site: Frontend (`/frontend`)

3. **Set Environment Variables**
   - Backend: `DATABASE_URL`, `PORT`, etc.
   - Frontend: `REACT_APP_API_URL`

#### Using Droplet (VPS)

1. **Create Droplet**
   - Choose Docker marketplace image
   - Select size (minimum 2GB RAM)

2. **SSH and Deploy**
   ```bash
   ssh root@your-droplet-ip
   git clone <your-repo>
   cd interview-prep-app
   docker-compose up -d
   ```

## Environment Variables

### Backend (.env)
```env
DATABASE_URL=postgresql://user:pass@host:port/dbname
PORT=8080
NODE_ENV=production
GIN_MODE=release
```

### Frontend (.env)
```env
REACT_APP_API_URL=https://your-backend-url/api/v1
```

## SSL/HTTPS Setup

### Using Certbot (VPS)
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### Using Cloudflare (Recommended)
1. Add your domain to Cloudflare
2. Update nameservers
3. Enable "Full SSL/TLS"
4. Enable "Always Use HTTPS"

## Monitoring

### Health Checks
- Backend: `GET /health`
- Frontend: `GET /health`

### Logging
```bash
# View logs
docker-compose logs -f

# View specific service
docker-compose logs -f backend
```

### Monitoring Services
- [UptimeRobot](https://uptimerobot.com) - Free uptime monitoring
- [Sentry](https://sentry.io) - Error tracking
- [LogDNA](https://logdna.com) - Log management

## Backup Strategy

### Database Backup
```bash
# Manual backup
docker exec interview-prep-db pg_dump -U interview_user interview_prep > backup.sql

# Automated backup (add to cron)
0 2 * * * docker exec interview-prep-db pg_dump -U interview_user interview_prep > /backups/backup_$(date +\%Y\%m\%d).sql
```

### Restore Backup
```bash
docker exec -i interview-prep-db psql -U interview_user interview_prep < backup.sql
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check DATABASE_URL format
   - Ensure database is running
   - Check network connectivity

2. **Frontend Can't Connect to Backend**
   - Verify REACT_APP_API_URL
   - Check CORS settings
   - Ensure backend is accessible

3. **Port Already in Use**
   ```bash
   # Find process using port
   lsof -i :8080
   # Kill process
   kill -9 <PID>
   ```

### Debug Commands
```bash
# Check container status
docker-compose ps

# View logs
docker-compose logs

# Access container shell
docker exec -it <container-name> sh

# Test database connection
docker exec -it interview-prep-db psql -U interview_user -d interview_prep
```

## CI/CD Setup

### GitHub Actions
Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Deploy to Railway
        uses: bervProject/railway-deploy@main
        with:
          railway_token: ${{ secrets.RAILWAY_TOKEN }}
```

## Security Best Practices

1. **Environment Variables**
   - Never commit .env files
   - Use strong passwords
   - Rotate credentials regularly

2. **Database**
   - Enable SSL/TLS
   - Restrict access by IP
   - Regular backups

3. **API**
   - Implement rate limiting
   - Add authentication (future)
   - Use HTTPS only

4. **Docker**
   - Keep images updated
   - Don't run as root
   - Scan for vulnerabilities

## Performance Optimization

1. **Frontend**
   - Enable gzip compression
   - Use CDN for static assets
   - Implement caching headers

2. **Backend**
   - Add database indexes
   - Implement caching layer
   - Use connection pooling

3. **Docker**
   - Multi-stage builds
   - Optimize image size
   - Use Alpine Linux

## Scaling

### Horizontal Scaling
```yaml
# docker-compose.yml
services:
  backend:
    deploy:
      replicas: 3
```

### Load Balancing
- Use Nginx or HAProxy
- Cloud provider load balancers
- Railway/Render auto-scaling

Remember to test your deployment thoroughly before going to production! 
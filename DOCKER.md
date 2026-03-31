# Docker Deployment Guide

This guide covers running the Agent Todo Platform with Docker.

## Prerequisites

- Docker Desktop (or Docker Engine + Docker Compose)
- At least 2GB of available RAM
- Ports 3000, 8080, and 5432 available

## Quick Start

### Development

```bash
# Start all services
make up
# or
./start.sh
# or
docker compose up -d
```

Access the application at:
- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **Swagger Docs:** http://localhost:8080/docs

### Production

```bash
# Create .env file with production values
cp .env.example .env
# Edit .env with your secure values

# Start production stack
make prod-up
# or
docker compose -f docker-compose.prod.yml up -d
```

## Available Commands

| Command | Description |
|---------|-------------|
| `make up` | Start all services in background |
| `make down` | Stop all services |
| `make build` | Build all Docker images |
| `make rebuild` | Rebuild all images without cache |
| `make logs` | View logs from all services |
| `make logs-b` | View backend logs |
| `make logs-f` | View frontend logs |
| `make ps` | Show running containers |
| `make clean` | Remove all containers, volumes, and images |

Or use the shell scripts:
- `./start.sh` - Start services
- `./stop.sh` - Stop services
- `./rebuild.sh` - Rebuild and restart

## Services

### Frontend (Port 3000)
- Vue 3 + Vite application
- Served by Nginx
- Proxies `/api/*` requests to backend

### Backend (Port 8080)
- Go + Gin application
- Connects to PostgreSQL
- Swagger docs at `/docs`

### Database (Port 5432)
- PostgreSQL 16
- Persistent volume for data

## Environment Variables

### Backend
| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://agent_todo:agent_todo_pass@postgres:5432/agent_todo?sslmode=disable` |
| `PORT` | Server port | `8080` |
| `GIN_MODE` | Gin mode (debug/release) | `release` |
| `JWT_SECRET` | JWT signing secret | (change in production!) |
| `FRONTEND_URL` | Frontend URL for CORS | `http://localhost:3000` |

### Database
| Variable | Description | Default |
|----------|-------------|---------|
| `POSTGRES_USER` | Database user | `agent_todo` |
| `POSTGRES_PASSWORD` | Database password | `agent_todo_pass` |
| `POSTGRES_DB` | Database name | `agent_todo` |

## Production Deployment

### 1. Create Environment File

```bash
cp .env.example .env
```

Edit `.env` with secure values:

```env
DB_PASSWORD=your-secure-password-here
JWT_SECRET=$(openssl rand -base64 32)
FRONTEND_URL=https://your-domain.com
```

### 2. Deploy

```bash
docker compose -f docker-compose.prod.yml up -d
```

### 3. Configure Reverse Proxy

Use a reverse proxy (nginx, traefik, caddy) to:
- Handle SSL/TLS termination
- Route traffic to containers
- Add security headers

Example Nginx config:

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using the port
lsof -i :3000
lsof -i :8080
lsof -i :5432

# Stop the conflicting service or change ports in docker-compose.yml
```

### Database Connection Issues

```bash
# Check database logs
make logs-db

# Connect to database
make shell-db
```

### View Container Logs

```bash
# All services
make logs

# Specific service
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f postgres
```

### Reset Everything

```bash
# Remove all containers, volumes, and local images
make clean

# Rebuild from scratch
make rebuild
```

## Health Checks

All services have health checks configured:

```bash
# Backend health
curl http://localhost:8080/health

# Check container health status
docker compose ps
```

## Getting Started

**First-time setup:** Create your account via the web UI at http://localhost:3000

⚠️ **For production:** Always use strong, unique passwords!

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Frontend      │────▶│    Backend      │────▶│   PostgreSQL    │
│   (Nginx)       │     │    (Go/Gin)     │     │                 │
│   Port 3000     │     │   Port 8080     │     │   Port 5432     │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

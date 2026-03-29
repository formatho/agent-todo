# Comprehensive Deployment Documentation & Runbooks

This document provides complete deployment procedures, runbooks, and operational guidance for the Agent Todo Platform.

## Table of Contents

- [1. System Architecture](#1-system-architecture)
- [2. Prerequisites](#2-prerequisites)
- [3. Production Server Setup](#3-production-server-setup)
- [4. First-time Deployment](#4-first-time-deployment)
- [5. Regular Deployment Process](#5-regular-deployment-process)
- [6. CI/CD Pipeline](#6-cicd-pipeline)
- [7. Monitoring & Health Checks](#7-monitoring--health-checks)
- [8. Security Best Practices](#8-security-best-practices)
- [9. Backup & Recovery](#9-backup--recovery)
- [10. Scaling & Performance](#10-scaling--performance)
- [11. Troubleshooting](#11-troubleshooting)
- [12. Maintenance Procedures](#12-maintenance-procedures)

## 1. System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Load Balancer / Reverse Proxy             │
│                   (Nginx / Traefik / Caddy)                   │
│                     (Port 443 - HTTPS)                       │
└─────────────────────────┬─────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      Production Server                        │
│                   (Droplet / VPS / AWS EC2)                   │
│                  /home/deploy/todo                           │
│                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   Frontend     │  │    Backend      │  │   PostgreSQL    │ │
│  │   (Nginx)      │  │    (Go/Gin)     │  │     (DB)        │ │
│  │  Port 3000     │  │  Port 18765     │  │   Port 5432     │ │
│  │                │  │                │  │                 │ │
│  │  Image:        │  │  Image:         │  │  Image:         │ │
│  │  ghcr.io/...   │  │  ghcr.io/...    │  │  postgres:16-   │ │
│  │  /frontend     │  │  /backend       │  │  alpine          │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 2. Prerequisites

### Server Requirements
- **OS**: Ubuntu 22.04 LTS or similar
- **CPU**: 2+ cores recommended
- **RAM**: 4GB+ (2GB for app, 2GB for database)
- **Storage**: 50GB+ SSD
- **Network**: Public IP with domain name
- **Ports**: 443 (HTTPS), 22 (SSH)

### Software Requirements
- Docker >= 20.10
- Docker Compose >= 2.0
- Nginx (for reverse proxy)
- Certbot (for SSL certificates)
- Git

### Domain & DNS
- Domain registered (e.g., `todo.formatho.com`)
- A record pointing to server IP
- Optional: Wildcard SSL certificate

## 3. Production Server Setup

### 3.1 Initial Server Setup

```bash
# Connect to server
ssh root@your-server-ip

# Update system
apt update && apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
usermod -aG docker $USER

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Install Nginx
apt install nginx -y

# Install Certbot
snap install core
snap refresh core
snap install --classic certbot
ln -s /snap/bin/certbot /usr/bin/certbot

# Create deploy user
adduser deploy
usermod -aG docker deploy

# Set up firewall
ufw allow OpenSSH
ufw allow 'Nginx Full'
ufw enable

# Configure SSH for deploy user
mkdir -p /home/deploy/.ssh
chmod 700 /home/deploy/.ssh
# Add SSH public key to /home/deploy/.ssh/authorized_keys
chmod 600 /home/deploy/.ssh/authorized_keys
chown -R deploy:deploy /home/deploy/.ssh

# Create application directory
mkdir -p /home/deploy/todo
chown deploy:deploy /home/deploy/todo
```

### 3.2 Network Configuration

```bash
# Create Docker network
docker network create formatho_app-network
```

### 3.3 Nginx Configuration

```bash
# Create nginx configuration
cat > /etc/nginx/sites-available/agent-todo << 'EOF'
server {
    listen 443 ssl http2;
    server_name todo.formatho.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/todo.formatho.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/todo.formatho.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;

    # Security headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req zone=api burst=20 nodelay;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:18765;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # CORS headers
        add_header Access-Control-Allow-Origin https://todo.formatho.com;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, PATCH, OPTIONS";
        add_header Access-Control-Allow-Headers "Origin, X-Requested-With, Content-Type, Accept, Authorization";
        
        # API rate limiting
        limit_req zone=api burst=30 nodelay;
    }

    # Health checks
    location /health {
        access_log off;
        proxy_pass http://localhost:18765/health;
    }

    # Static files caching
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        proxy_pass http://localhost:3000;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}

# HTTP redirect to HTTPS
server {
    listen 80;
    server_name todo.formatho.com;
    return 301 https://$host$request_uri;
}
EOF

# Enable site
ln -s /etc/nginx/sites-available/agent-todo /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Test nginx configuration
nginx -t

# Reload nginx
systemctl reload nginx
```

### 3.4 SSL Certificate Setup

```bash
# Obtain SSL certificate
certbot --nginx -d todo.formatho.com

# Set up auto-renewal
echo "0 12 * * * /usr/bin/certbot renew --quiet" | crontab -
```

## 4. First-time Deployment

### 4.1 Repository Setup

```bash
# Switch to deploy user
su - deploy

# Clone repository
cd /home/deploy
git clone https://github.com/formatho/agent-todo.git todo
cd todo

# Create necessary directories
mkdir -p logs
mkdir -p data/postgres
chmod 755 logs data/postgres

# Configure environment
cat > .env << 'EOF'
# Database
DB_USER=agent_todo
DB_NAME=agent_todo
DB_PASSWORD=your-secure-database-password-here

# Backend
JWT_SECRET=your-secure-jwt-secret-here
FRONTEND_URL=https://todo.formatho.com

# Deployment
COMPOSE_PROJECT_NAME=agent-todo
EOF

# Set secure file permissions
chmod 600 .env
```

### 4.2 Docker Configuration

```bash
# Create production docker-compose override
cat > docker-compose.override.prod.yml << 'EOF'
version: '3.8'

services:
  postgres:
    environment:
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - ./data/postgres:/var/lib/postgresql/data
    networks:
      - formatho_app-network

  backend:
    environment:
      DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable
      PORT: 18765
      GIN_MODE: release
      JWT_SECRET: ${JWT_SECRET}
      FRONTEND_URL: ${FRONTEND_URL}
      DISABLE_REGISTRATION: "true"
    restart: unless-stopped
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - formatho_app-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:18765/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  frontend:
    restart: unless-stopped
    depends_on:
      - backend
    networks:
      - formatho_app-network

networks:
  formatho_app-network:
    external: true
EOF
```

### 4.3 Initial Deployment

```bash
# Pull latest images
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml pull

# Start services
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml up -d

# Check status
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml ps

# Verify health
curl -f https://todo.formatho.com/health

# View logs
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml logs -f backend
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml logs -f frontend
```

## 5. Regular Deployment Process

### 5.1 Manual Deployment

```bash
# SSH to server and switch to deploy user
ssh deploy@your-server-ip
cd /home/deploy/todo

# Pull latest changes
git pull origin main

# Rebuild and deploy
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml pull
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml up -d

# Verify deployment
curl -f https://todo.formatho.com/health
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml ps

# Check logs if needed
docker compose logs --tail=50 backend
```

### 5.2 Rolling Deployment Strategy

```bash
# Update specific service
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml pull backend
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml up -d backend

# Wait for health check
sleep 10
curl -f https://todo.formatho.com/api/health || exit 1

# Update frontend
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml pull frontend
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml up -d frontend

# Final verification
curl -f https://todo.formatho.com/health
```

## 6. CI/CD Pipeline

### 6.1 GitHub Secrets Configuration

Required secrets in GitHub repository:
- `TODO_DEPLOY_HOST`: Server IP/hostname
- `TODO_DEPLOY_USER`: deploy username
- `TODO_DEPLOY_KEY`: SSH private key
- `TODO_GHCR_TOKEN`: GitHub Container Registry token

### 6.2 Pipeline Flow

1. **Push to main branch**
   - Backend: Lint → Test → Build → Push → Deploy
   - Frontend: Lint → Build → Push → Deploy

2. **Automatic Health Checks**
   - Deployment waits for service health
   - Rollback if health checks fail

3. **Notification System**
   - Success/failure notifications
   - Deployment logs available

### 6.3 Pipeline Verification

```bash
# Check recent deployments
gh run list --limit 5

# View deployment logs
gh run view <run-id> --job deploy

# Check deployed versions
docker ps | grep agent-todo
```

## 7. Monitoring & Health Checks

### 7.1 Service Health Checks

```bash
# Backend health
curl -f https://todo.formatho.com/api/health

# Database connectivity
docker exec agent-todo-db pg_isready -U agent_todo

# Container status
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml ps

# Resource usage
docker stats --no-stream
```

### 7.2 Monitoring Setup

```bash
# Create monitoring scripts
cat > /home/deploy/todo/health-check.sh << 'EOF'
#!/bin/bash

# Health check script
set -e

echo "🔍 Running health checks..."

# Backend health
if ! curl -f https://todo.formatho.com/api/health > /dev/null 2>&1; then
    echo "❌ Backend health check failed"
    exit 1
fi

# Database health
if ! docker exec agent-todo-db pg_isready -U agent_todo > /dev/null 2>&1; then
    echo "❌ Database health check failed"
    exit 1
fi

echo "✅ All health checks passed"
EOF

chmod +x /home/deploy/todo/health-check.sh
```

### 7.3 Automated Monitoring

```bash
# Add to crontab
echo "*/5 * * * * /home/deploy/todo/health-check.sh" | crontab -

# Alert configuration (optional)
# Integrate with monitoring tools like:
# - UptimeRobot
# - Pingdom
# - Custom webhook alerts
```

## 8. Security Best Practices

### 8.1 Container Security

```bash
# Regular security updates
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml pull

# Scan images for vulnerabilities
docker scan ghcr.io/formatho/agent-todo/backend:latest
docker scan ghcr.io/formatho/agent-todo/frontend:latest

# Use specific tags instead of 'latest' in production
# Update docker-compose.prod.yml to use versioned tags
```

### 8.2 Network Security

```bash
# Firewall rules
ufw status

# Allow only necessary ports
ufw allow 22/tcp   # SSH
ufw allow 443/tcp  # HTTPS
ufw deny 3000/tcp  # Block direct access to frontend
ufw deny 18765/tcp # Block direct access to backend
```

### 8.3 Application Security

```bash
# Regular dependency updates
cd /home/deploy/todo/backend
go mod tidy
go mod download

# Check for vulnerabilities
go list -u -m all
```

## 9. Backup & Recovery

### 9.1 Database Backup

```bash
# Create backup script
cat > /home/deploy/todo/backup-db.sh << 'EOF'
#!/bin/bash

# Database backup script
set -e

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/home/deploy/todo/backups"
DB_NAME="agent_todo"

mkdir -p "$BACKUP_DIR"

# Create backup
docker exec agent-todo-db pg_dump -U agent_todo agent_todo > "$BACKUP_DIR/db_backup_$DATE.sql"

# Compress backup
gzip "$BACKUP_DIR/db_backup_$DATE.sql"

# Keep only last 7 days of backups
find "$BACKUP_DIR" -name "db_backup_*.sql.gz" -mtime +7 -delete

echo "✅ Database backup completed: db_backup_$DATE.sql.gz"
EOF

chmod +x /home/deploy/todo/backup-db.sh

# Add to crontab for daily backups
echo "0 2 * * * /home/deploy/todo/backup-db.sh" | crontab -
```

### 9.2 File Backup

```bash
# Configuration backup
tar -czf /home/deploy/todo/backups/config_backup_$(date +%Y%m%d).tar.gz \
  --exclude=node_modules \
  --exclude=.git \
  --exclude=logs \
  --exclude=data/postgres \
  /home/deploy/todo/
```

### 9.3 Disaster Recovery

```bash
# Recovery procedure
# 1. Restore database from backup
docker exec -i agent-todo-db psql -U agent_todo agent_todo < /path/to/backup.sql

# 2. Restore configuration files
tar -xzf /path/to/config_backup.tar.gz -C /home/deploy/todo/

# 3. Restart services
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml restart
```

## 10. Scaling & Performance

### 10.1 Horizontal Scaling

```bash
# Scale backend services
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml up -d --scale backend=3

# Load balancing with nginx
# Update nginx config to include upstream backend
cat > /etc/nginx/conf.d/upstream.conf << 'EOF'
upstream backend {
    server localhost:18765;
    server localhost:18766;
    server localhost:18767;
}

server {
    location /api/ {
        proxy_pass http://backend;
        # ... rest of config
    }
}
EOF
```

### 10.2 Resource Optimization

```bash
# Monitor resource usage
docker stats

# Optimize database
docker exec agent-todo-db psql -U agent_todo -c "VACUUM ANALYZE;"

# Clean up unused containers and images
docker system prune -f
```

## 11. Troubleshooting

### 11.1 Common Issues

**502 Bad Gateway**
```bash
# Check backend logs
docker compose logs backend

# Verify backend health
curl http://localhost:18765/health

# Check nginx configuration
nginx -t
```

**Database Connection Issues**
```bash
# Check database logs
docker compose logs postgres

# Test database connectivity
docker exec -it agent-todo-db psql -U agent_todo -c "SELECT 1;"

# Reset database connection
docker compose restart postgres
```

**SSL Certificate Issues**
```bash
# Check certificate expiration
certbot certificates

# Renew certificate
certbot renew --dry-run

# Test nginx configuration
nginx -t
systemctl reload nginx
```

### 11.2 Debug Mode

```bash
# Enable debug logging
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml up -d \
  -e DEBUG=true backend

# View detailed logs
docker compose logs -f backend --tail=100
```

## 12. Maintenance Procedures

### 12.1 Regular Maintenance Schedule

```bash
# Weekly maintenance
# Sundays at 3 AM
cat > /home/deploy/todo/maintenance-weekly.sh << 'EOF'
#!/bin/bash

echo "🔧 Running weekly maintenance..."

# Update system packages (if needed)
# sudo apt update && sudo apt upgrade -y

# Clean up Docker
docker system prune -f
docker volume prune -f

# Rotate logs
docker compose logs --tail=100 > /home/deploy/todo/logs/$(date +%Y%m%d)_rotation.log

# Health check
/home/deploy/todo/health-check.sh

echo "✅ Weekly maintenance completed"
EOF

chmod +x /home/deploy/todo/maintenance-weekly.sh
echo "0 3 * * 0 /home/deploy/todo/maintenance-weekly.sh" | crontab -
```

### 12.2 Version Updates

```bash
# Update procedure
cat > /home/deploy/todo/update-version.sh << 'EOF'
#!/bin/bash

# Update application version
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

echo "🔄 Updating to version $VERSION..."

# Stop services
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml down

# Update docker-compose.prod.yml to use specific version
sed -i "s/:latest/:$VERSION/g" docker-compose.prod.yml

# Pull new versions
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml pull

# Start services
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml up -d

# Health check
sleep 10
/home/deploy/todo/health-check.sh

echo "✅ Update to version $VERSION completed"
EOF

chmod +x /home/deploy/todo/update-version.sh
```

### 12.3 Emergency Procedures

**Emergency Stop**
```bash
# Emergency stop all services
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml down

# Keep database running
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml up -d postgres
```

**Emergency Restart**
```bash
# Force restart all services
docker compose -f docker-compose.prod.yml -f docker-compose.override.prod.yml restart

# Wait for health checks
sleep 30
/home/deploy/todo/health-check.sh
```

---

## Quick Reference Commands

### Development
```bash
# Start development environment
make up

# View logs
make logs

# Stop services
make down
```

### Production
```bash
# Deploy
cd /home/deploy/todo && git pull && docker compose up -d

# Health checks
curl -f https://todo.formatho.com/health

# View logs
docker compose logs -f backend

# Maintenance
./maintenance-weekly.sh
```

### Emergency
```bash
# Emergency stop
docker compose down

# Restart everything
docker compose up -d

# Check status
docker compose ps
```

---

*Last updated: March 28, 2026*
*Version: 1.0*
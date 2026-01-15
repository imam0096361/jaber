# Docker Compose Quick Reference

## Basic Commands

### Start Services
```bash
# Start all services in background
docker compose up -d

# Start with development profile (includes Adminer)
docker compose --profile dev up -d

# Start with production profile (includes Nginx)
docker compose --profile production up -d

# Start and rebuild
docker compose up -d --build
```

### Stop Services
```bash
# Stop all services
docker compose down

# Stop and remove volumes (WARNING: deletes data!)
docker compose down -v

# Stop specific service
docker compose stop app
```

### View Logs
```bash
# View all logs
docker compose logs

# Follow logs in real-time
docker compose logs -f

# View specific service logs
docker compose logs app
docker compose logs postgresdb

# Last 100 lines
docker compose logs --tail=100 app
```

### Container Management
```bash
# List running containers
docker compose ps

# Restart service
docker compose restart app

# Execute command in container
docker compose exec app sh
docker compose exec postgresdb psql -U postgres -d news

# View container stats
docker stats
```

## Database Operations

### Backup Database
```bash
# Create backup
docker compose exec postgresdb pg_dump -U postgres news > backup_$(date +%Y%m%d).sql

# Create compressed backup
docker compose exec postgresdb pg_dump -U postgres news | gzip > backup_$(date +%Y%m%d).sql.gz
```

### Restore Database
```bash
# Restore from backup
docker compose exec -T postgresdb psql -U postgres news < backup.sql

# Restore from compressed backup
gunzip < backup.sql.gz | docker compose exec -T postgresdb psql -U postgres news
```

### Access Database
```bash
# Using psql
docker compose exec postgresdb psql -U postgres -d news

# Using Adminer (if enabled)
# Open browser: http://localhost:8080
```

## Maintenance

### Update Application
```bash
# Pull latest code
git pull

# Rebuild and restart
docker compose down
docker compose up -d --build

# Rolling update (zero downtime)
docker compose up -d --no-deps --build app
```

### Clean Up
```bash
# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune

# Remove everything unused
docker system prune -a --volumes

# View disk usage
docker system df
```

### View Resource Usage
```bash
# Real-time stats
docker stats

# Container resource limits
docker compose config
```

## Troubleshooting

### View Container Details
```bash
# Inspect container
docker compose exec app env

# Check health status
docker compose ps

# View container processes
docker compose top
```

### Debug Issues
```bash
# Check logs for errors
docker compose logs app | grep -i error

# Access container shell
docker compose exec app sh

# Check network
docker network ls
docker network inspect news_news-network
```

### Restart Services
```bash
# Restart all
docker compose restart

# Restart specific service
docker compose restart app

# Force recreate
docker compose up -d --force-recreate app
```

## Environment Profiles

### Development
```bash
docker compose --profile dev up -d
# Includes: PostgreSQL + App + Adminer
```

### Production
```bash
docker compose up -d
# Includes: PostgreSQL + App
```

### Production with Nginx
```bash
docker compose --profile production up -d
# Includes: PostgreSQL + App + Nginx
```

## Health Checks

### Manual Health Check
```bash
# Check app health
curl http://localhost:9998/v1/health-check

# Check database
docker compose exec postgresdb pg_isready -U postgres
```

### View Health Status
```bash
docker compose ps
# Look for "healthy" status
```

## Port Mappings

- **9998**: Application (host) → 3000 (container)
- **5432**: PostgreSQL (host) → 5432 (container)
- **8080**: Adminer (host) → 8080 (container) [dev profile]
- **80/443**: Nginx (host) → 80/443 (container) [production profile]

## Volume Management

### List Volumes
```bash
docker volume ls | grep news
```

### Backup Volumes
```bash
# Backup postgres data
docker run --rm -v news_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz /data
```

### Restore Volumes
```bash
# Restore postgres data
docker run --rm -v news_postgres_data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres_backup.tar.gz -C /
```

## Security

### Update Images
```bash
# Pull latest images
docker compose pull

# Restart with new images
docker compose up -d
```

### View Container Logs for Security
```bash
# Check for failed login attempts
docker compose logs app | grep -i "failed\|unauthorized"
```

## Quick Deployment

### First Time Setup
```bash
# 1. Configure environment
cp .env.production .env
nano .env

# 2. Start services
docker compose up -d

# 3. Check status
docker compose ps
docker compose logs -f
```

### Update Deployment
```bash
# Pull, rebuild, and restart
git pull
docker compose down
docker compose up -d --build
```

## Monitoring

### Watch Logs
```bash
# Follow all logs
docker compose logs -f

# Follow specific service
docker compose logs -f app

# Filter logs
docker compose logs app | grep ERROR
```

### Resource Monitoring
```bash
# Real-time stats
docker stats

# Disk usage
docker system df
```

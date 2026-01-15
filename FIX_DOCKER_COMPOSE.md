# Quick Fix Guide for Docker Compose Error

## Problem
You're encountering a `KeyError: 'ContainerConfig'` error when running `make docker`. This is a known bug in older versions of docker-compose (v1.29.2).

## Solution Options

### Option 1: Use the Fix Script (Quickest)
```bash
chmod +x fix-docker-compose.sh
./fix-docker-compose.sh
```

### Option 2: Manual Steps
```bash
# Stop and remove everything
docker-compose down -v --remove-orphans

# Clean up Docker system
docker system prune -f

# Start fresh
docker-compose up -d --build
```

### Option 3: Upgrade Docker Compose (Recommended for long-term)
```bash
# Remove old docker-compose
sudo apt-get remove docker-compose

# Install latest docker-compose v2 (comes with Docker)
sudo apt-get update
sudo apt-get install docker-compose-plugin

# Verify installation
docker compose version
```

**Note:** Docker Compose v2 uses `docker compose` (space) instead of `docker-compose` (hyphen).

## After Fixing

Check if containers are running:
```bash
docker-compose ps
# or with v2:
docker compose ps
```

View logs:
```bash
docker-compose logs -f app
```

## Fix Upload Permissions
After containers are running, fix the upload permissions:
```bash
chmod +x fix-upload-permission.sh
./fix-upload-permission.sh
```

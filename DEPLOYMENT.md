# 🚀 Ubuntu Server Deployment Guide

Complete guide to deploy the News Application on Ubuntu Server using Docker.

---

## 📋 Prerequisites

### System Requirements
- **Ubuntu Server**: 20.04 LTS or newer
- **RAM**: Minimum 2GB (4GB recommended)
- **Disk Space**: Minimum 10GB free
- **CPU**: 2 cores recommended

### Required Software
- Docker Engine
- Docker Compose
- Git (optional, for cloning)

---

## 🔧 Step 1: Install Docker on Ubuntu Server

### Update System
```bash
sudo apt update
sudo apt upgrade -y
```

### Install Docker
```bash
# Install prerequisites
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common

# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Add Docker repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Verify installation
docker --version
docker compose version
```

### Configure Docker (Optional but Recommended)
```bash
# Add your user to docker group (to run docker without sudo)
sudo usermod -aG docker $USER

# Apply group changes (or logout/login)
newgrp docker

# Enable Docker to start on boot
sudo systemctl enable docker
sudo systemctl start docker
```

---

## 📦 Step 2: Upload Application Files

### Option A: Using Git (Recommended)
```bash
# Clone the repository
git clone <your-repo-url> /opt/news-app
cd /opt/news-app
```

### Option B: Using SCP/SFTP
```bash
# From your local machine (Windows)
# Using PowerShell or Command Prompt
scp -r C:\Users\ihcking\Downloads\news username@your-server-ip:/opt/news-app

# Or use FileZilla, WinSCP, or similar tools
```

### Option C: Manual Upload
1. Use an FTP client (FileZilla, WinSCP)
2. Upload all files to `/opt/news-app` on the server

---

## ⚙️ Step 3: Configure Environment

### Create Production Environment File
```bash
cd /opt/news-app

# Copy the production template
cp .env.production .env

# Edit the environment file
nano .env
```

### Update These Critical Values:
```bash
# IMPORTANT: Change these values!
DB_PASSWORD=your_very_secure_database_password_here
JWT_SECRET=your_very_long_random_secret_key_at_least_32_characters_long

# Optional: Update these if needed
APP_PORT=9998
DB_NAME=news
DB_USER=postgres
```

**Security Tips:**
- Use strong, random passwords (minimum 16 characters)
- Generate JWT secret: `openssl rand -base64 32`
- Never commit `.env` to version control

---

## 🐳 Step 4: Build and Run with Docker

### Basic Deployment (Development/Testing)
```bash
cd /opt/news-app

# Build and start all services
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f
```

### Production Deployment (with Nginx)
```bash
# Start with production profile
docker compose --profile production up -d

# This starts: PostgreSQL + App + Nginx
```

### With Database Admin Panel (Development)
```bash
# Start with admin profile
docker compose --profile dev up -d

# Access Adminer at: http://your-server-ip:8080
```

---

## 🔍 Step 5: Verify Deployment

### Check Container Status
```bash
# View running containers
docker compose ps

# Should show:
# - news-postgres (healthy)
# - news-app (healthy)
# - news-adminer (optional)
# - news-nginx (optional)
```

### Check Application Logs
```bash
# View all logs
docker compose logs

# View specific service logs
docker compose logs app
docker compose logs postgresdb

# Follow logs in real-time
docker compose logs -f app
```

### Test Application
```bash
# Test from server
curl http://localhost:9998/

# Test health check
curl http://localhost:9998/v1/health-check
```

---

## 🌐 Step 6: Access the Application

### From Browser
- **Home**: `http://your-server-ip:9998/`
- **Login**: `http://your-server-ip:9998/login`
- **Admin**: `http://your-server-ip:9998/admin`
- **Adminer** (if enabled): `http://your-server-ip:8080/`

### Configure Firewall
```bash
# Allow HTTP traffic
sudo ufw allow 9998/tcp

# Allow HTTPS (if using Nginx with SSL)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Allow Adminer (optional, for development only)
sudo ufw allow 8080/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

---

## 🗄️ Step 7: Database Management

### Run Migrations
```bash
# Access the app container
docker compose exec app sh

# Inside container, run migrations (if needed)
# Or use the migration commands from Makefile
```

### Backup Database
```bash
# Create backup
docker compose exec postgresdb pg_dump -U postgres news > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore from backup
docker compose exec -T postgresdb psql -U postgres news < backup_file.sql
```

### Access Database Directly
```bash
# Using docker exec
docker compose exec postgresdb psql -U postgres -d news

# Or use Adminer web interface (if enabled)
# http://your-server-ip:8080
```

---

## 🔄 Step 8: Update and Maintenance

### Update Application
```bash
cd /opt/news-app

# Pull latest changes (if using Git)
git pull

# Rebuild and restart
docker compose down
docker compose up -d --build

# Or rolling update (zero downtime)
docker compose up -d --no-deps --build app
```

### View Resource Usage
```bash
# Container stats
docker stats

# Disk usage
docker system df
```

### Clean Up
```bash
# Remove stopped containers
docker compose down

# Remove with volumes (WARNING: deletes data!)
docker compose down -v

# Clean up unused images
docker system prune -a
```

---

## 🛡️ Step 9: Security Hardening (Production)

### 1. Use HTTPS with SSL/TLS
```bash
# Install Certbot for Let's Encrypt
sudo apt install -y certbot

# Generate SSL certificate
sudo certbot certonly --standalone -d your-domain.com

# Copy certificates to nginx/ssl/
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem nginx/ssl/key.pem

# Update nginx.conf to enable HTTPS
# Restart nginx
docker compose restart nginx
```

### 2. Disable Adminer in Production
```bash
# Remove --profile dev when starting
docker compose up -d  # Without profiles
```

### 3. Regular Updates
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Update Docker images
docker compose pull
docker compose up -d
```

### 4. Monitor Logs
```bash
# Set up log rotation
sudo nano /etc/docker/daemon.json
```

Add:
```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

```bash
# Restart Docker
sudo systemctl restart docker
```

---

## 📊 Step 10: Monitoring (Optional)

### View Application Logs
```bash
# Real-time logs
docker compose logs -f app

# Last 100 lines
docker compose logs --tail=100 app
```

### Health Checks
```bash
# Check container health
docker compose ps

# Manual health check
curl http://localhost:9998/v1/health-check
```

### Set Up Monitoring (Advanced)
Consider using:
- **Prometheus + Grafana** for metrics
- **ELK Stack** for log aggregation
- **Uptime Kuma** for uptime monitoring

---

## 🚨 Troubleshooting

### Container Won't Start
```bash
# Check logs
docker compose logs app

# Check if port is in use
sudo netstat -tulpn | grep 9998

# Restart services
docker compose restart
```

### Database Connection Issues
```bash
# Check database is running
docker compose ps postgresdb

# Check database logs
docker compose logs postgresdb

# Verify credentials in .env file
cat .env | grep DB_
```

### Permission Issues
```bash
# Fix uploads directory permissions
sudo chown -R 1000:1000 frontend/uploads

# Restart app
docker compose restart app
```

### Out of Disk Space
```bash
# Clean up Docker
docker system prune -a --volumes

# Check disk usage
df -h
docker system df
```

---

## 📝 Useful Commands Reference

```bash
# Start services
docker compose up -d

# Stop services
docker compose down

# Restart specific service
docker compose restart app

# View logs
docker compose logs -f

# Execute command in container
docker compose exec app sh

# Rebuild and restart
docker compose up -d --build

# Scale services (if needed)
docker compose up -d --scale app=3

# Export/Import database
docker compose exec postgresdb pg_dump -U postgres news > backup.sql
docker compose exec -T postgresdb psql -U postgres news < backup.sql
```

---

## 🎯 Quick Start Commands

```bash
# Complete deployment in one go
cd /opt/news-app
cp .env.production .env
nano .env  # Update passwords and secrets
docker compose up -d
docker compose logs -f
```

---

## 📞 Support

If you encounter issues:
1. Check logs: `docker compose logs -f`
2. Verify environment: `cat .env`
3. Check container status: `docker compose ps`
4. Review this guide's troubleshooting section

---

## ✅ Deployment Checklist

- [ ] Ubuntu Server updated
- [ ] Docker and Docker Compose installed
- [ ] Application files uploaded
- [ ] `.env` file configured with secure passwords
- [ ] Firewall configured
- [ ] Application started with `docker compose up -d`
- [ ] Health check passed
- [ ] Application accessible from browser
- [ ] Database backup strategy in place
- [ ] SSL/TLS configured (production)
- [ ] Monitoring set up (optional)

---

**🎉 Congratulations! Your application is now deployed!**

Access it at: `http://your-server-ip:9998`

# 🚀 News Application - Docker Deployment Package

## 📦 What's Included

This package contains everything you need to deploy the News Application on an Ubuntu server using Docker.

### ✅ Files Created/Updated

1. **Dockerfile** - Optimized multi-stage production build
2. **docker-compose.yml** - Complete orchestration with profiles
3. **.env.production** - Production environment template
4. **.dockerignore** - Optimized build context
5. **nginx/nginx.conf** - Reverse proxy configuration
6. **deploy.sh** - Automated deployment script
7. **DEPLOYMENT.md** - Comprehensive deployment guide
8. **DOCKER_COMMANDS.md** - Quick command reference

---

## 🎯 Quick Start (3 Steps)

### On Ubuntu Server:

```bash
# 1. Upload files to server (or clone from git)
cd /opt/news-app

# 2. Run the deployment script
chmod +x deploy.sh
./deploy.sh

# 3. Access your application
# http://your-server-ip:9998
```

That's it! The script handles everything automatically.

---

## 📋 Manual Deployment (Alternative)

If you prefer manual control:

```bash
# 1. Configure environment
cp .env.production .env
nano .env  # Update DB_PASSWORD and JWT_SECRET

# 2. Start services
docker compose up -d

# 3. Check status
docker compose ps
docker compose logs -f
```

---

## 🔧 Key Features

### ✨ Production-Ready
- Multi-stage Docker build (smaller image size)
- Non-root user for security
- Health checks for all services
- Automatic restarts
- Log rotation support

### 🛡️ Security
- Secure default configurations
- Environment-based secrets
- Optional SSL/TLS with Nginx
- Rate limiting
- Isolated network

### 📊 Flexible Deployment
- **Development**: With Adminer for database management
- **Production**: Minimal services
- **Production + Nginx**: With reverse proxy and SSL

### 🔄 Easy Maintenance
- One-command updates
- Database backup/restore scripts
- Health monitoring
- Resource management

---

## 🌐 Deployment Modes

### Mode 1: Development
```bash
docker compose --profile dev up -d
```
**Includes**: PostgreSQL + App + Adminer
**Use for**: Testing, development, debugging

### Mode 2: Production
```bash
docker compose up -d
```
**Includes**: PostgreSQL + App
**Use for**: Simple production deployment

### Mode 3: Production + Nginx
```bash
docker compose --profile production up -d
```
**Includes**: PostgreSQL + App + Nginx
**Use for**: Production with reverse proxy, SSL, load balancing

---

## 📁 Project Structure

```
news/
├── Dockerfile                 # Application container definition
├── docker-compose.yml         # Service orchestration
├── .env                       # Environment configuration (create from .env.production)
├── .env.production           # Production template
├── .dockerignore             # Build optimization
├── deploy.sh                 # Automated deployment script
├── DEPLOYMENT.md             # Full deployment guide
├── DOCKER_COMMANDS.md        # Command reference
├── nginx/
│   └── nginx.conf            # Reverse proxy config
├── frontend/                 # Frontend files
│   ├── index.html
│   ├── login.html
│   ├── admin.html
│   ├── dashboard.html
│   ├── css/
│   ├── js/
│   └── uploads/              # User uploads (persistent)
├── src/                      # Go backend source
│   ├── main.go
│   ├── config/
│   ├── database/
│   │   ├── migrations/
│   │   └── init/
│   ├── middleware/
│   ├── router/
│   └── utils/
├── go.mod
└── go.sum
```

---

## 🔑 Environment Variables

### Required (Must Change in Production!)
```bash
DB_PASSWORD=your_secure_password_here
JWT_SECRET=your_very_long_random_secret_key
```

### Optional (Can Keep Defaults)
```bash
APP_ENV=production
APP_PORT=9998
DB_HOST=postgresdb
DB_PORT=5432
DB_USER=postgres
DB_NAME=news
JWT_ACCESS_EXP_MINUTES=60
JWT_REFRESH_EXP_DAYS=7
```

---

## 🚀 Deployment Checklist

Before deploying to production:

- [ ] Ubuntu Server 20.04+ ready
- [ ] Docker and Docker Compose installed
- [ ] Files uploaded to server
- [ ] `.env` file created and configured
- [ ] **DB_PASSWORD** changed from default
- [ ] **JWT_SECRET** set to random string (32+ chars)
- [ ] Firewall configured (ports 9998, 80, 443)
- [ ] Domain name configured (if using)
- [ ] SSL certificate obtained (if using HTTPS)
- [ ] Backup strategy planned

---

## 📊 Access Points

After deployment, access:

| Service | URL | Notes |
|---------|-----|-------|
| **Application** | `http://server-ip:9998` | Main application |
| **Login** | `http://server-ip:9998/login` | User login |
| **Admin** | `http://server-ip:9998/admin` | Admin panel |
| **Adminer** | `http://server-ip:8080` | DB admin (dev mode only) |
| **Health Check** | `http://server-ip:9998/v1/health-check` | API health |

---

## 🛠️ Common Commands

```bash
# Start services
docker compose up -d

# Stop services
docker compose down

# View logs
docker compose logs -f

# Restart app
docker compose restart app

# Update and redeploy
git pull && docker compose up -d --build

# Backup database
docker compose exec postgresdb pg_dump -U postgres news > backup.sql

# Check status
docker compose ps
```

See **DOCKER_COMMANDS.md** for complete reference.

---

## 📖 Documentation

1. **DEPLOYMENT.md** - Complete step-by-step deployment guide
   - Docker installation
   - Server setup
   - Security hardening
   - Troubleshooting

2. **DOCKER_COMMANDS.md** - Quick command reference
   - Common operations
   - Database management
   - Maintenance tasks

3. **README.md** - Original project documentation

---

## 🔒 Security Best Practices

### ✅ Implemented
- Multi-stage builds (minimal attack surface)
- Non-root container user
- Environment-based secrets
- Health checks
- Network isolation
- Resource limits

### 📝 Recommended
1. Change default passwords in `.env`
2. Use strong JWT secret (32+ characters)
3. Enable SSL/TLS in production
4. Disable Adminer in production
5. Set up firewall rules
6. Regular security updates
7. Database backups
8. Log monitoring

---

## 🆘 Troubleshooting

### Application won't start
```bash
# Check logs
docker compose logs app

# Check environment
docker compose exec app env | grep DB_
```

### Database connection failed
```bash
# Check database is running
docker compose ps postgresdb

# Check credentials
cat .env | grep DB_
```

### Port already in use
```bash
# Check what's using the port
sudo netstat -tulpn | grep 9998

# Change port in .env
APP_PORT=8080
```

See **DEPLOYMENT.md** for detailed troubleshooting.

---

## 📞 Support Resources

- **Deployment Guide**: See `DEPLOYMENT.md`
- **Command Reference**: See `DOCKER_COMMANDS.md`
- **Docker Logs**: `docker compose logs -f`
- **Container Status**: `docker compose ps`

---

## 🎉 Next Steps

1. **Deploy**: Run `./deploy.sh` on your Ubuntu server
2. **Configure**: Update `.env` with secure credentials
3. **Test**: Access the application in your browser
4. **Secure**: Set up SSL/TLS for production
5. **Monitor**: Set up logging and monitoring
6. **Backup**: Implement database backup strategy

---

## 📝 Notes

- **Development**: Use `--profile dev` for Adminer access
- **Production**: Use default profile (no Adminer)
- **SSL/TLS**: Use `--profile production` with Nginx
- **Updates**: Simply run `docker compose up -d --build`
- **Backups**: Automated backup scripts included

---

## ✨ Features

### Application Features
- User authentication (JWT)
- Admin panel
- File uploads
- RESTful API
- Swagger documentation

### Infrastructure Features
- Docker containerization
- PostgreSQL database
- Nginx reverse proxy (optional)
- Health monitoring
- Auto-restart on failure
- Persistent data volumes
- Log management

---

## 🔄 Update Process

```bash
# Pull latest changes
git pull

# Rebuild and restart
docker compose down
docker compose up -d --build

# Verify
docker compose ps
docker compose logs -f app
```

---

**Ready to deploy? Start with `./deploy.sh` or follow `DEPLOYMENT.md`!** 🚀

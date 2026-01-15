# News Application

A modern news management application built with Go (Fiber) backend and PostgreSQL database.

## 🚀 Quick Start

### Using Docker (Recommended)
```bash
# Start the application
docker compose up -d

# Access the application
# http://localhost:2345
```

### Using Make
```bash
# Start the application
make start

# Access the application
# http://localhost:2345
```

## 📋 Access Points

- **Home**: http://localhost:2345/
- **Login**: http://localhost:2345/login
- **Admin**: http://localhost:2345/admin
- **API Health**: http://localhost:2345/v1/health-check

## 🐳 Docker Deployment

### Quick Deploy on Ubuntu Server
```bash
# Run the automated deployment script
chmod +x deploy.sh
./deploy.sh
```

### Manual Deployment
```bash
# 1. Configure environment
cp .env.production .env
nano .env  # Update DB_PASSWORD and JWT_SECRET

# 2. Start services
docker compose up -d

# 3. Check status
docker compose ps
```

## 📚 Documentation

- **[DOCKER_README.md](DOCKER_README.md)** - Docker deployment overview
- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Complete deployment guide
- **[DOCKER_COMMANDS.md](DOCKER_COMMANDS.md)** - Command reference
- **[RUN.md](RUN.md)** - Quick run instructions

## 🔧 Technology Stack

- **Backend**: Go (Fiber framework)
- **Database**: PostgreSQL
- **Frontend**: HTML, CSS, JavaScript
- **Containerization**: Docker & Docker Compose
- **Reverse Proxy**: Nginx (optional)

## ⚙️ Configuration

### Environment Variables

Key variables in `.env`:
```bash
APP_PORT=2345
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=root
DB_NAME=news
JWT_SECRET=change-me
```

## 🛠️ Development

### Prerequisites
- Go 1.22+
- PostgreSQL 16+
- Docker & Docker Compose (for containerized deployment)

### Run Locally
```bash
# Install dependencies
go mod download

# Run migrations
make migrate-up

# Start server
make start
```

### Run Tests
```bash
make tests
```

## 🐳 Docker Deployment Modes

### Development (with Adminer)
```bash
docker compose --profile dev up -d
```
Access Adminer at: http://localhost:8080

### Production
```bash
docker compose up -d
```

### Production with Nginx
```bash
docker compose --profile production up -d
```

## 📦 Project Structure

```
news/
├── frontend/          # Frontend files
├── src/               # Go backend source
├── nginx/             # Nginx configuration
├── Dockerfile         # Docker build file
├── docker-compose.yml # Docker orchestration
├── Makefile           # Build commands
└── deploy.sh          # Deployment script
```

## 🔒 Security

- Change `DB_PASSWORD` in production
- Set a strong `JWT_SECRET` (32+ characters)
- Use HTTPS in production
- Disable Adminer in production

## 📝 License

MIT License

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## 📞 Support

For deployment help, see [DEPLOYMENT.md](DEPLOYMENT.md)

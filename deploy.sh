#!/bin/bash

# News Application - Quick Deploy Script for Ubuntu Server
# This script automates the deployment process

set -e  # Exit on error

echo "=================================================="
echo "  News Application - Ubuntu Server Deployment"
echo "=================================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    print_warning "Please do not run this script as root"
    exit 1
fi

# Step 1: Check Docker installation
echo "Step 1: Checking Docker installation..."
if command -v docker &> /dev/null; then
    print_success "Docker is installed: $(docker --version)"
else
    print_error "Docker is not installed"
    echo "Would you like to install Docker? (y/n)"
    read -r install_docker
    if [ "$install_docker" = "y" ]; then
        echo "Installing Docker..."
        sudo apt update
        sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
        sudo apt update
        sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
        sudo usermod -aG docker $USER
        print_success "Docker installed successfully"
        print_warning "Please log out and log back in for group changes to take effect"
        exit 0
    else
        print_error "Docker is required. Exiting."
        exit 1
    fi
fi

# Step 2: Check Docker Compose
echo ""
echo "Step 2: Checking Docker Compose..."
if docker compose version &> /dev/null; then
    print_success "Docker Compose is installed: $(docker compose version)"
else
    print_error "Docker Compose is not installed"
    exit 1
fi

# Step 3: Check environment file
echo ""
echo "Step 3: Checking environment configuration..."
if [ ! -f ".env" ]; then
    if [ -f ".env.production" ]; then
        print_warning ".env file not found. Copying from .env.production"
        cp .env.production .env
        print_info "Please edit .env file and update the following:"
        print_info "  - DB_PASSWORD (set a strong password)"
        print_info "  - JWT_SECRET (set a long random string)"
        echo ""
        echo "Edit .env now? (y/n)"
        read -r edit_env
        if [ "$edit_env" = "y" ]; then
            ${EDITOR:-nano} .env
        else
            print_warning "Remember to edit .env before running in production!"
        fi
    else
        print_error ".env file not found and no template available"
        exit 1
    fi
else
    print_success ".env file found"
    
    # Check for default passwords
    if grep -q "CHANGE_THIS" .env; then
        print_warning "Default passwords detected in .env file!"
        print_warning "Please update DB_PASSWORD and JWT_SECRET"
        echo "Edit .env now? (y/n)"
        read -r edit_env
        if [ "$edit_env" = "y" ]; then
            ${EDITOR:-nano} .env
        fi
    fi
fi

# Step 4: Create necessary directories
echo ""
echo "Step 4: Creating necessary directories..."
mkdir -p frontend/uploads
mkdir -p nginx/ssl
print_success "Directories created"

# Step 5: Choose deployment mode
echo ""
echo "Step 5: Choose deployment mode:"
echo "  1) Development (with Adminer)"
echo "  2) Production (without Adminer)"
echo "  3) Production with Nginx reverse proxy"
read -p "Enter choice (1-3): " deploy_mode

COMPOSE_PROFILES=""
case $deploy_mode in
    1)
        COMPOSE_PROFILES="--profile dev"
        print_info "Deploying in development mode with Adminer"
        ;;
    2)
        print_info "Deploying in production mode"
        ;;
    3)
        COMPOSE_PROFILES="--profile production"
        print_info "Deploying with Nginx reverse proxy"
        print_warning "Make sure nginx/nginx.conf is properly configured"
        ;;
    *)
        print_error "Invalid choice"
        exit 1
        ;;
esac

# Step 6: Pull latest images
echo ""
echo "Step 6: Pulling Docker images..."
docker compose pull
print_success "Images pulled"

# Step 7: Build and start services
echo ""
echo "Step 7: Building and starting services..."
docker compose $COMPOSE_PROFILES up -d --build

# Wait for services to be healthy
echo ""
echo "Step 8: Waiting for services to be healthy..."
sleep 5

# Check container status
echo ""
echo "Container Status:"
docker compose ps

# Step 9: Display access information
echo ""
echo "=================================================="
echo "  Deployment Complete!"
echo "=================================================="
echo ""

# Get server IP
SERVER_IP=$(hostname -I | awk '{print $1}')

print_success "Application is running!"
echo ""
echo "Access URLs:"
echo "  🏠 Home:     http://$SERVER_IP:9998/"
echo "  🔐 Login:    http://$SERVER_IP:9998/login"
echo "  👨‍💼 Admin:    http://$SERVER_IP:9998/admin"

if [ "$deploy_mode" = "1" ]; then
    echo "  🗄️  Adminer:  http://$SERVER_IP:8080/"
fi

if [ "$deploy_mode" = "3" ]; then
    echo "  🌐 Nginx:    http://$SERVER_IP/"
fi

echo ""
echo "Useful commands:"
echo "  View logs:        docker compose logs -f"
echo "  Stop services:    docker compose down"
echo "  Restart:          docker compose restart"
echo "  View status:      docker compose ps"
echo ""

# Step 10: Check health
echo "Checking application health..."
sleep 10

if curl -f http://localhost:9998/v1/health-check &> /dev/null; then
    print_success "Health check passed!"
else
    print_warning "Health check failed. Check logs with: docker compose logs app"
fi

echo ""
print_info "For detailed documentation, see DEPLOYMENT.md"
echo ""

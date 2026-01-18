#!/bin/bash

# =====================================================
# Fix Database Password Authentication Error
# =====================================================
# This script fixes the SASL authentication error by
# resetting the PostgreSQL password to match the .env file
# =====================================================

set -e

echo "🔧 Fixing Database Password Authentication..."
echo ""

# Navigate to project directory (update this path for your server)
cd /home/star/jaber 2>/dev/null || cd /home/star 2>/dev/null || echo "Using current directory"

# Read the password from .env file
if [ -f .env ]; then
    DB_PASSWORD=$(grep -E "^DB_PASSWORD=" .env | cut -d'=' -f2)
    DB_USER=$(grep -E "^DB_USER=" .env | cut -d'=' -f2)
    DB_NAME=$(grep -E "^DB_NAME=" .env | cut -d'=' -f2)
else
    DB_PASSWORD="root"
    DB_USER="postgres"
    DB_NAME="news"
fi

echo "📋 Current configuration:"
echo "   DB_USER: $DB_USER"
echo "   DB_PASSWORD: $DB_PASSWORD"
echo "   DB_NAME: $DB_NAME"
echo ""

# Check if PostgreSQL container is running
if ! docker ps | grep -q "news-postgres"; then
    echo "❌ PostgreSQL container is not running!"
    echo "   Starting containers..."
    docker-compose up -d postgresdb
    sleep 10
fi

echo "🔑 Resetting PostgreSQL password..."

# Method 1: Try to change password using ALTER USER
docker exec -it news-postgres psql -U postgres -c "ALTER USER postgres WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null && {
    echo "✅ Password updated successfully!"
} || {
    echo "⚠️  Method 1 failed, trying Method 2..."
    
    # Method 2: If that fails, we need to recreate the database with correct password
    echo ""
    echo "🗑️  This requires recreating the database container."
    echo "   ⚠️  WARNING: This will delete all existing data!"
    echo ""
    read -p "Do you want to proceed? (y/N): " confirm
    
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        echo "Stopping containers..."
        docker-compose down
        
        echo "Removing old postgres volume..."
        docker volume rm jaber_postgres_data 2>/dev/null || docker volume rm $(docker volume ls -q | grep postgres) 2>/dev/null || true
        
        echo "Recreating containers with correct password..."
        docker-compose up -d --build
        
        echo "Waiting for services to start..."
        sleep 15
        
        echo "✅ Database recreated with correct password!"
    else
        echo "❌ Aborted. Please fix manually."
        exit 1
    fi
}

echo ""
echo "🔄 Restarting app container to apply changes..."
docker-compose restart app

echo ""
echo "⏳ Waiting for app to start..."
sleep 10

echo ""
echo "🧪 Testing database connection..."
docker exec news-postgres psql -U postgres -d news -c "SELECT COUNT(*) as article_count FROM articles;" 2>/dev/null && {
    echo ""
    echo "✅ Database connection successful!"
} || {
    echo "⚠️  Table might not exist yet (this is OK for fresh database)"
}

echo ""
echo "🎉 Fix complete! Try creating an article now."
echo ""
echo "📝 If you still have issues, run:"
echo "   docker-compose logs -f app"

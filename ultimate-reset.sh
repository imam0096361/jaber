#!/bin/bash

# =====================================================
# ULTIMATE DATABASE RESET & FIX SCRIPT
# =====================================================
# This script permanently fixes all database connection
# issues by:
# 1. Backing up existing data
# 2. Removing old volumes with potentially wrong password
# 3. Recreating everything with correct, consistent credentials
# =====================================================

set -e

echo "╔════════════════════════════════════════════════════════╗"
echo "║        ULTIMATE DATABASE RESET & FIX SCRIPT            ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""

# Navigate to project directory
cd /home/star/jaber 2>/dev/null || {
    echo "❌ Could not find /home/star/jaber"
    echo "   Please run this script from the jaber directory"
    exit 1
}

echo "📁 Working directory: $(pwd)"
echo ""

# Fixed credentials (hardcoded for consistency)
DB_USER="postgres"
DB_PASSWORD="root"
DB_NAME="news"

echo "📋 Database Configuration:"
echo "   User: $DB_USER"
echo "   Password: $DB_PASSWORD"
echo "   Database: $DB_NAME"
echo ""

# Step 1: Backup existing data if possible
echo "═══════════════════════════════════════════════════════"
echo "STEP 1: Attempting to backup existing data..."
echo "═══════════════════════════════════════════════════════"

BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"

if docker ps | grep -q "news-postgres"; then
    echo "📦 Backing up existing data..."
    docker exec news-postgres pg_dump -U $DB_USER $DB_NAME > "$BACKUP_FILE" 2>/dev/null && {
        echo "✅ Backup saved to: $BACKUP_FILE"
    } || {
        echo "⚠️  Could not backup (database may be empty or inaccessible)"
    }
else
    echo "⚠️  PostgreSQL container not running, skipping backup"
fi

echo ""

# Step 2: Stop all containers
echo "═══════════════════════════════════════════════════════"
echo "STEP 2: Stopping all containers..."
echo "═══════════════════════════════════════════════════════"

docker-compose down 2>/dev/null || docker compose down 2>/dev/null || {
    echo "⚠️  docker-compose not found, stopping manually..."
    docker stop news-app news-postgres news-adminer 2>/dev/null || true
    docker rm news-app news-postgres news-adminer 2>/dev/null || true
}

echo "✅ Containers stopped"
echo ""

# Step 3: Remove old database volume
echo "═══════════════════════════════════════════════════════"
echo "STEP 3: Removing old database volume..."
echo "═══════════════════════════════════════════════════════"

# Find and remove postgres volumes
VOLUMES=$(docker volume ls -q | grep -E "(postgres|jaber)" 2>/dev/null || true)
if [ -n "$VOLUMES" ]; then
    echo "Found volumes: $VOLUMES"
    for vol in $VOLUMES; do
        docker volume rm "$vol" 2>/dev/null && echo "✅ Removed: $vol" || echo "⚠️  Could not remove: $vol"
    done
else
    echo "⚠️  No postgres volumes found"
fi

# Also try the specific name patterns
docker volume rm jaber_postgres_data 2>/dev/null || true
docker volume rm news_postgres_data 2>/dev/null || true
docker volume rm postgres_data 2>/dev/null || true

echo ""

# Step 4: Pull latest code
echo "═══════════════════════════════════════════════════════"
echo "STEP 4: Pulling latest code..."
echo "═══════════════════════════════════════════════════════"

git pull origin main 2>/dev/null || git pull 2>/dev/null || echo "⚠️  Git pull skipped"
echo ""

# Step 5: Rebuild and start containers
echo "═══════════════════════════════════════════════════════"
echo "STEP 5: Rebuilding and starting containers..."
echo "═══════════════════════════════════════════════════════"

docker-compose up -d --build 2>/dev/null || docker compose up -d --build

echo "⏳ Waiting for containers to start (30 seconds)..."
sleep 30

echo ""

# Step 6: Verify database connection
echo "═══════════════════════════════════════════════════════"
echo "STEP 6: Verifying database connection..."
echo "═══════════════════════════════════════════════════════"

MAX_ATTEMPTS=10
for i in $(seq 1 $MAX_ATTEMPTS); do
    echo "Attempt $i/$MAX_ATTEMPTS..."
    if docker exec news-postgres psql -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
        echo "✅ Database connection successful!"
        break
    fi
    if [ $i -eq $MAX_ATTEMPTS ]; then
        echo "❌ Database connection failed after $MAX_ATTEMPTS attempts"
        echo "   Check logs: docker logs news-postgres"
        exit 1
    fi
    sleep 5
done

echo ""

# Step 7: Verify app connection
echo "═══════════════════════════════════════════════════════"
echo "STEP 7: Verifying app connection to database..."
echo "═══════════════════════════════════════════════════════"

sleep 10

if docker logs news-app 2>&1 | grep -q "Database connection established"; then
    echo "✅ App connected to database successfully!"
else
    echo "Checking app logs..."
    docker logs news-app --tail 20
fi

echo ""

# Step 8: Restore backup if exists
echo "═══════════════════════════════════════════════════════"
echo "STEP 8: Checking for backup restoration..."
echo "═══════════════════════════════════════════════════════"

if [ -f "$BACKUP_FILE" ] && [ -s "$BACKUP_FILE" ]; then
    echo "📦 Found backup file: $BACKUP_FILE"
    read -p "Do you want to restore the backup? (y/N): " restore
    if [ "$restore" = "y" ] || [ "$restore" = "Y" ]; then
        docker exec -i news-postgres psql -U $DB_USER -d $DB_NAME < "$BACKUP_FILE" && {
            echo "✅ Backup restored successfully!"
        } || {
            echo "⚠️  Backup restoration had some errors (this is often OK)"
        }
    else
        echo "Skipping backup restoration"
    fi
else
    echo "No backup file to restore"
fi

echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║                    ✅ COMPLETE!                        ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo "📊 Container Status:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "(news|NAMES)"
echo ""
echo "🌐 Your app should now be running at: http://103.118.19.134:2345"
echo ""
echo "📝 If you still have issues, check logs with:"
echo "   docker logs news-app -f"
echo "   docker logs news-postgres -f"
echo ""

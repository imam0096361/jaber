#!/bin/bash

echo "=========================================="
echo "  ULTIMATE RESET - Complete Fresh Start"
echo "=========================================="
echo ""
echo "⚠️  WARNING: This will DELETE ALL DATA!"
echo ""

cd /home/star/jaber/jaber

echo "Step 1: Stopping all containers..."
docker-compose down

echo ""
echo "Step 2: Removing all volumes (deleting database)..."
docker-compose down -v --remove-orphans

echo ""
echo "Step 3: Removing old images..."
docker image rm jaber_app 2>/dev/null || true
docker image prune -f

echo ""
echo "Step 4: Cleaning up Docker system..."
docker system prune -f

echo ""
echo "Step 5: Setting correct permissions..."
chmod -R 755 ./src/database/init 2>/dev/null || true
mkdir -p ./frontend/uploads
chmod -R 777 ./frontend/uploads

echo ""
echo "Step 6: Rebuilding and starting containers..."
docker-compose up -d --build

echo ""
echo "Step 7: Waiting for containers to initialize..."
sleep 15

echo ""
echo "Step 8: Checking container status..."
docker-compose ps

echo ""
echo "Step 9: Checking database connection..."
docker exec news-postgres psql -U postgres -d news -c "SELECT COUNT(*) FROM articles;" 2>/dev/null && echo "✅ Database connected!" || echo "⏳ Database still initializing..."

echo ""
echo "=========================================="
echo "  ✅ RESET COMPLETE!"
echo "=========================================="
echo ""
echo "Access points:"
echo "  - Home:  http://$(hostname -I | awk '{print $1}'):2345/"
echo "  - Admin: http://$(hostname -I | awk '{print $1}'):2345/admin"
echo "  - Login: http://$(hostname -I | awk '{print $1}'):2345/login"
echo ""
echo "View logs: docker-compose logs -f app"
echo ""

#!/bin/bash

# Quick fix script for database connection issue
# Run this on your Ubuntu server

echo "=================================================="
echo "  Fixing Database Connection Issue"
echo "=================================================="
echo ""

cd /home/star/jaber/jaber

echo "Step 1: Pulling latest fixes from GitHub..."
git pull

echo ""
echo "Step 2: Stopping containers..."
docker compose down

echo ""
echo "Step 3: Rebuilding with fixes..."
docker compose up -d --build

echo ""
echo "Step 4: Waiting for services to start..."
sleep 10

echo ""
echo "Step 5: Checking container status..."
docker compose ps

echo ""
echo "Step 6: Checking application logs..."
docker compose logs --tail=20 app

echo ""
echo "=================================================="
echo "  Fix Applied!"
echo "=================================================="
echo ""
echo "Check if the app is running:"
echo "  docker compose ps"
echo ""
echo "View logs:"
echo "  docker compose logs -f app"
echo ""
echo "Test the application:"
echo "  curl http://localhost:2345/v1/health-check"
echo ""

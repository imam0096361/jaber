#!/bin/bash

echo "==================================================="
echo "  Restarting Docker Containers"
echo "==================================================="
echo ""

echo "Step 1: Stopping containers..."
docker-compose down

echo ""
echo "Step 2: Starting containers in detached mode..."
docker-compose up -d --build

echo ""
echo "Step 3: Waiting for containers to start..."
sleep 5

echo ""
echo "Step 4: Checking container status..."
docker-compose ps

echo ""
echo "==================================================="
echo "  ✓ Containers Restarted!"
echo "==================================================="
echo ""
echo "View logs with: docker-compose logs -f app"
echo "Check status with: docker-compose ps"
echo ""

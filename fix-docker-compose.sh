#!/bin/bash

# Fix Docker Compose ContainerConfig error
# This script stops and removes old containers, then starts fresh

echo "==================================================="
echo "  Fixing Docker Compose ContainerConfig Error"
echo "==================================================="
echo ""

echo "Step 1: Stopping all containers..."
docker-compose down

echo ""
echo "Step 2: Removing old containers and volumes..."
docker-compose down -v --remove-orphans

echo ""
echo "Step 3: Pruning Docker system (optional but recommended)..."
docker system prune -f

echo ""
echo "Step 4: Starting fresh containers..."
docker-compose up -d --build

echo ""
echo "==================================================="
echo "  ✓ Docker Compose Fixed!"
echo "==================================================="
echo ""
echo "Check status with: docker-compose ps"
echo "View logs with: docker-compose logs -f"
echo ""

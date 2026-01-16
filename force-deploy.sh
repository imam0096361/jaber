#!/bin/bash

echo "=========================================="
echo "  FORCE CACHE CLEAR - Deploy New Version"
echo "=========================================="

cd /home/star/jaber/jaber

echo ""
echo "Step 1: Pull latest code..."
git pull

echo ""
echo "Step 2: Stop containers..."
docker-compose down

echo ""
echo "Step 3: Remove old images to force rebuild..."
docker image prune -f

echo ""
echo "Step 4: Rebuild and start containers..."
docker-compose up -d --build --force-recreate

echo ""
echo "Step 5: Wait for containers..."
sleep 5

echo ""
echo "Step 6: Show container status..."
docker-compose ps

echo ""
echo "=========================================="
echo "  ✅ DONE! New version deployed"
echo "=========================================="
echo ""
echo "NOW ON YOUR COMPUTER:"
echo "1. Close ALL browser tabs of the site"
echo "2. Clear browser cache (Ctrl+Shift+Delete)"  
echo "3. Close and reopen browser completely"
echo "4. Open site in NEW INCOGNITO window"
echo ""
echo "OR just visit: http://103.118.19.134:2345/?nocache=$(date +%s)"
echo ""

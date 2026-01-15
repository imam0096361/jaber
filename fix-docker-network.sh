#!/bin/bash

# Quick Fix for Docker Network Issues
# Run this script on your Ubuntu server to fix IPv6 connectivity problems

echo "=================================================="
echo "  Docker Network Fix - Disabling IPv6"
echo "=================================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root or with sudo"
    echo "Usage: sudo bash fix-docker-network.sh"
    exit 1
fi

echo "Step 1: Creating Docker daemon configuration..."

# Create docker directory if it doesn't exist
mkdir -p /etc/docker

# Create daemon.json with IPv6 disabled
cat > /etc/docker/daemon.json <<EOF
{
  "ipv6": false,
  "dns": ["8.8.8.8", "8.8.4.4"],
  "ip6tables": false
}
EOF

echo "✓ Docker daemon configuration created"
echo ""

echo "Step 2: Restarting Docker service..."
systemctl restart docker

# Wait for Docker to restart
sleep 3

echo "✓ Docker service restarted"
echo ""

echo "Step 3: Verifying Docker status..."
if systemctl is-active --quiet docker; then
    echo "✓ Docker is running"
else
    echo "✗ Docker failed to start"
    echo "Check logs with: journalctl -u docker -n 50"
    exit 1
fi

echo ""
echo "Step 4: Testing Docker connectivity..."
if docker pull hello-world > /dev/null 2>&1; then
    echo "✓ Docker can connect to Docker Hub"
    docker rmi hello-world > /dev/null 2>&1
else
    echo "✗ Still having connectivity issues"
    echo "Try running: docker pull hello-world"
    exit 1
fi

echo ""
echo "=================================================="
echo "  ✓ Fix Applied Successfully!"
echo "=================================================="
echo ""
echo "You can now deploy your application:"
echo "  cd /home/star/jaber/jaber"
echo "  git pull"
echo "  docker compose up -d --build"
echo ""

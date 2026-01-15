# 🔧 Network Connectivity Fix for Docker

## Problem
Your Ubuntu server is experiencing IPv6 connectivity issues when trying to reach Docker Hub:
```
network is unreachable: dial tcp [2600:1f18:2148:bc00:6511:89ca:d5b7:769c]:443
```

## Solutions

### Solution 1: Disable IPv6 for Docker (Recommended)

This forces Docker to use IPv4 only:

```bash
# Create or edit Docker daemon configuration
sudo nano /etc/docker/daemon.json
```

Add the following content:
```json
{
  "ipv6": false,
  "ip6tables": false
}
```

Then restart Docker:
```bash
sudo systemctl restart docker
```

### Solution 2: Configure DNS to Use IPv4

```bash
# Edit Docker daemon config
sudo nano /etc/docker/daemon.json
```

Add:
```json
{
  "dns": ["8.8.8.8", "8.8.4.4"],
  "ipv6": false
}
```

Restart Docker:
```bash
sudo systemctl restart docker
```

### Solution 3: Disable IPv6 System-Wide

```bash
# Edit sysctl configuration
sudo nano /etc/sysctl.conf
```

Add these lines:
```
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
net.ipv6.conf.lo.disable_ipv6 = 1
```

Apply changes:
```bash
sudo sysctl -p
```

Restart Docker:
```bash
sudo systemctl restart docker
```

### Solution 4: Use Docker Mirror (Alternative)

If you're in a region with Docker Hub access issues, use a mirror:

```bash
sudo nano /etc/docker/daemon.json
```

Add:
```json
{
  "registry-mirrors": ["https://mirror.gcr.io"],
  "ipv6": false
}
```

Restart Docker:
```bash
sudo systemctl restart docker
```

## Quick Fix (Try This First)

Run these commands on your Ubuntu server:

```bash
# 1. Create Docker daemon config
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json > /dev/null <<EOF
{
  "ipv6": false,
  "dns": ["8.8.8.8", "8.8.4.4"]
}
EOF

# 2. Restart Docker
sudo systemctl restart docker

# 3. Verify Docker is running
sudo systemctl status docker

# 4. Test connectivity
docker pull hello-world

# 5. Now try deploying again
cd /home/star/jaber/jaber
git pull
docker compose up -d --build
```

## Verify Fix

After applying the fix, test Docker connectivity:

```bash
# Test Docker
docker run hello-world

# Should see: "Hello from Docker!"
```

## Alternative: Pre-built Images

If network issues persist, you can build the Docker image on a machine with good connectivity and transfer it:

### On a machine with good internet:
```bash
# Build the image
docker build -t news-app:latest .

# Save to file
docker save news-app:latest > news-app.tar

# Transfer to server (using scp)
scp news-app.tar user@your-server:/tmp/
```

### On your Ubuntu server:
```bash
# Load the image
docker load < /tmp/news-app.tar

# Update docker-compose.yml to use the loaded image
# Change: build: .
# To: image: news-app:latest

# Start services
docker compose up -d
```

## After Fix

Once network is fixed, deploy normally:

```bash
cd /home/star/jaber/jaber
docker compose down
docker compose up -d --build
```

## Check Status

```bash
# View logs
docker compose logs -f

# Check containers
docker compose ps

# Test application
curl http://localhost:2345/v1/health-check
```

---

**Recommended**: Use **Solution 1** (disable IPv6 for Docker) - it's the quickest and most reliable fix.

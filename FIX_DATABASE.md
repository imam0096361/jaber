# 🔧 Database Connection Fix

## Problem
The application was trying to connect to `localhost` instead of the PostgreSQL container `postgresdb`.

**Error seen:**
```
failed to connect to `host=localhost user=postgres database=news`: 
dial error (dial tcp [::1]:5432: connect: connection refused)
```

## Root Cause
The Dockerfile was copying the `.env` file which had `DB_HOST=localhost`. This overrode the environment variables set by docker-compose.yml.

## Solution Applied
✅ Removed `.env` file copy from Dockerfile  
✅ Docker Compose now provides all environment variables  
✅ `DB_HOST=postgresdb` is correctly set in docker-compose.yml

---

## 🚀 Quick Fix (Run on Your Server)

### Option 1: Automated Script
```bash
cd /home/star/jaber/jaber
git pull
bash fix-db-connection.sh
```

### Option 2: Manual Steps
```bash
# Navigate to project
cd /home/star/jaber/jaber

# Pull latest fixes
git pull

# Stop containers
docker compose down

# Rebuild and start
docker compose up -d --build

# Check status
docker compose ps

# View logs
docker compose logs -f app
```

---

## ✅ Verification

After applying the fix, you should see:

### 1. Containers Running
```bash
docker compose ps
```
Expected output:
```
NAME            STATUS
news-app        Up (healthy)
news-postgres   Up (healthy)
```

### 2. No Database Errors in Logs
```bash
docker compose logs app | grep -i error
```
Should show no connection errors.

### 3. Application Accessible
```bash
curl http://localhost:2345/v1/health-check
```
Should return a successful response.

### 4. Database Connected
```bash
docker compose logs app | grep -i "database"
```
Should show successful database connection messages.

---

## 🔍 Understanding the Fix

### Before (Incorrect)
```dockerfile
# Dockerfile was copying .env with DB_HOST=localhost
COPY --from=builder /app/.env .
```

The app read `.env` file → `DB_HOST=localhost` → tried to connect to localhost → failed (PostgreSQL is in another container)

### After (Correct)
```dockerfile
# Dockerfile doesn't copy .env
# Environment variables come from docker-compose.yml
```

Docker Compose sets `DB_HOST=postgresdb` → app uses this → connects to PostgreSQL container → success!

---

## 📊 Environment Variables Flow

```
docker-compose.yml
    ↓
Sets environment variables:
    - DB_HOST=postgresdb
    - DB_PORT=5432
    - DB_USER=postgres
    - DB_PASSWORD=root
    - DB_NAME=news
    ↓
Container starts with these variables
    ↓
Application reads environment variables
    ↓
Connects to PostgreSQL container successfully
```

---

## 🛠️ Troubleshooting

### If still seeing connection errors:

#### 1. Check if PostgreSQL is running
```bash
docker compose ps postgresdb
```
Should show "Up (healthy)"

#### 2. Check PostgreSQL logs
```bash
docker compose logs postgresdb
```
Look for any errors

#### 3. Verify network
```bash
docker network ls | grep news
docker network inspect jaber_news-network
```
Both containers should be on the same network

#### 4. Test database connection manually
```bash
# From the app container
docker compose exec app sh
# Inside container, try to ping postgres
ping postgresdb
# Should resolve and respond
```

#### 5. Check environment variables in container
```bash
docker compose exec app env | grep DB_
```
Should show:
```
DB_HOST=postgresdb
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=root
DB_NAME=news
```

---

## 🔄 Complete Deployment Steps

For a clean deployment:

```bash
# 1. Navigate to project
cd /home/star/jaber/jaber

# 2. Pull latest code
git pull

# 3. Stop and remove old containers
docker compose down -v

# 4. Build and start fresh
docker compose up -d --build

# 5. Wait for services to be healthy
sleep 15

# 6. Check status
docker compose ps

# 7. View logs
docker compose logs -f

# 8. Test application
curl http://localhost:2345/v1/health-check
```

---

## 📝 What Changed

### Files Modified:
1. **Dockerfile** - Removed `.env` file copy
2. **fix-db-connection.sh** - Added automated fix script

### Files Unchanged (Already Correct):
- **docker-compose.yml** - Already had correct `DB_HOST=postgresdb`
- **All other configuration files**

---

## ✅ Expected Successful Output

After running the fix, you should see logs like:

```
news-app  | Config file loaded from ./
news-app  | Successfully connected to database
news-app  | Database migration completed
news-app  | 
news-app  |  ┌───────────────────────────────────────────────────┐
news-app  |  │                     Fiber API                     │
news-app  |  │                  Fiber v2.52.10                   │
news-app  |  │               http://0.0.0.0:3000                 │
news-app  |  │                                                   │
news-app  |  │ Handlers ........... 101  Processes ........... 1 │
news-app  |  │ Prefork ....... Disabled  PID ................. 1 │
news-app  |  └───────────────────────────────────────────────────┘
```

No error messages about database connection!

---

## 🎯 Quick Commands

```bash
# Pull and apply fix
cd /home/star/jaber/jaber && git pull && docker compose down && docker compose up -d --build

# Check status
docker compose ps

# View logs
docker compose logs -f app

# Test health
curl http://localhost:2345/v1/health-check

# Access application
curl http://localhost:2345/
```

---

## 📞 Still Having Issues?

If the problem persists:

1. **Check Docker logs**: `docker compose logs`
2. **Verify network**: `docker network inspect jaber_news-network`
3. **Check PostgreSQL**: `docker compose exec postgresdb psql -U postgres -d news -c "\l"`
4. **Restart everything**: `docker compose down -v && docker compose up -d --build`

---

**The fix is ready! Just run `git pull` and rebuild.** 🚀

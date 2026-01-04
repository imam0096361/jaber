# PostgreSQL Setup for News Portal

## Installation

### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

### macOS
```bash
brew install postgresql@15
brew services start postgresql@15
```

### Windows
Download and install from: https://www.postgresql.org/download/windows/

## Configuration

### 1. Start PostgreSQL Service
```bash
# Linux
sudo systemctl start postgresql

# macOS
brew services start postgresql@15

# Windows
# Usually starts automatically
```

### 2. Create Database and User
```bash
# Connect to PostgreSQL
sudo -u postgres psql

# Run these commands in PostgreSQL:
CREATE DATABASE news_portal;
CREATE USER news_user WITH ENCRYPTED PASSWORD 'news_password';
ALTER ROLE news_user SET client_encoding TO 'utf8';
ALTER ROLE news_user SET default_transaction_isolation TO 'read committed';
ALTER ROLE news_user SET default_transaction_deferrable TO on;
ALTER ROLE news_user SET default_transaction_read_committed TO on;
GRANT ALL PRIVILEGES ON DATABASE news_portal TO news_user;
\q
```

### 3. Update .env File
```env
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=news_user
DB_PASSWORD=news_password
DB_NAME=news_portal
```

### 4. Run the Application
```bash
go run ./src/main.go
```

## Database Tables

The application will automatically create these tables:

### articles
- id (SERIAL PRIMARY KEY)
- title (VARCHAR 255)
- content (TEXT)
- category (VARCHAR 100)
- author (VARCHAR 100)
- image (VARCHAR 255)
- created (TIMESTAMP)
- featured (BOOLEAN)

### users
- id (SERIAL PRIMARY KEY)
- name (VARCHAR 100)
- email (VARCHAR 100 UNIQUE)
- password_hash (VARCHAR 255)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

## Verification

```bash
# Connect to database
psql -U news_user -d news_portal

# List tables
\dt

# Exit
\q
```

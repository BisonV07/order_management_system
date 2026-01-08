# Setup Without Docker

If you don't have Docker or prefer to use a local PostgreSQL installation, follow these instructions.

## Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- PostgreSQL installed locally (version 12+)

## Step 1: Install PostgreSQL

### macOS
```bash
# Using Homebrew
brew install postgresql@15
brew services start postgresql@15
```

### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

### Windows
Download and install from: https://www.postgresql.org/download/windows/

## Step 2: Create Database

```bash
# Connect to PostgreSQL
psql postgres

# Create database and user
CREATE DATABASE oms_db;
CREATE USER oms_user WITH PASSWORD 'oms_password';
GRANT ALL PRIVILEGES ON DATABASE oms_db TO oms_user;

# Exit psql
\q
```

## Step 3: Setup Backend

### 3.1 Navigate to backend directory
```bash
cd backend
```

### 3.2 Install Go dependencies
```bash
go mod download
go mod tidy
```

### 3.3 Configure environment
```bash
cp .env.sample .env
```

Edit `.env` file with your local PostgreSQL credentials:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=oms_user
DB_PASSWORD=oms_password
DB_NAME=oms_db
DB_SSLMODE=disable
SERVER_PORT=8080
JWT_SECRET=your-secret-key-here
JWT_EXPIRY=24h
LOG_LEVEL=info
```

### 3.4 Run database migrations
```bash
go run cmd/main.go --migrate
```

### 3.5 Start the API server
```bash
go run cmd/main.go --api --port=8080
```

## Step 4: Setup Frontend

### 4.1 Navigate to frontend directory (in a new terminal)
```bash
cd frontend
```

### 4.2 Install dependencies
```bash
npm install
```

### 4.3 Configure environment
```bash
cp .env.sample .env
```

The `.env` file should contain:
```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 4.4 Start the development server
```bash
npm run dev
```

## Verify Setup

1. **Database**: Check if PostgreSQL is running
   ```bash
   # macOS/Linux
   psql -U oms_user -d oms_db -c "SELECT version();"
   
   # Or check if service is running
   # macOS
   brew services list | grep postgresql
   
   # Linux
   sudo systemctl status postgresql
   ```

2. **Backend API**: Test health endpoint
   ```bash
   curl http://localhost:8080/api/v1/health
   ```
   Should return: `OK`

3. **Frontend**: Open browser to `http://localhost:3000`

## Troubleshooting

### PostgreSQL Connection Issues

**Error: "connection refused"**
- Ensure PostgreSQL is running:
  ```bash
  # macOS
  brew services start postgresql@15
  
  # Linux
  sudo systemctl start postgresql
  ```

**Error: "password authentication failed"**
- Check your `.env` file has correct credentials
- Reset password if needed:
  ```bash
  psql postgres
  ALTER USER oms_user WITH PASSWORD 'oms_password';
  ```

**Error: "database does not exist"**
- Create the database:
  ```bash
  psql postgres
  CREATE DATABASE oms_db;
  ```

### Port Already in Use

If port 5432 is already in use:
1. Find what's using it:
   ```bash
   # macOS/Linux
   lsof -i :5432
   ```
2. Either stop that service or change PostgreSQL port in `backend/.env`

## Quick Start (All Commands)

```bash
# Terminal 1: Backend
cd backend
go mod download
cp .env.sample .env
# Edit .env with your PostgreSQL credentials
go run cmd/main.go --migrate
go run cmd/main.go --api --port=8080

# Terminal 2: Frontend
cd frontend
npm install
cp .env.sample .env
npm run dev
```


# Setup Instructions - Run These Commands

Since automated installation requires your password, please run these commands manually:

## Step 1: Install Docker Desktop

```bash
brew install --cask docker
```

**After installation completes:**
1. Open Docker Desktop from Applications folder
2. Wait for Docker to fully start (you'll see a whale icon in your menu bar)
3. Verify it's running: `docker ps` (should not show an error)

## Step 2: Start PostgreSQL Container

Once Docker is running, execute:

```bash
cd /Users/sapaharamsingh/intern_demo/order_management_system
docker compose up -d postgres
```

Wait a few seconds, then verify:
```bash
docker ps | grep postgres
```

## Step 3: Create Configuration Files

Run the setup script to create .env files:
```bash
./setup.sh
```

Or manually create them:
```bash
cd backend
# Create .env file with database settings (see backend/.env.sample for template)
go mod download
```

## Step 4: Run Database Migrations

```bash
go run cmd/main.go --migrate
```

## Step 5: Start Backend Server

```bash
go run cmd/main.go --api --port=8080
```

## Step 6: Setup Frontend (in a new terminal)

```bash
cd frontend
cp .env.sample .env
npm install
npm run dev
```

## Quick Copy-Paste Commands

Run these in order:

```bash
# 1. Install Docker (enter password when prompted)
brew install --cask docker

# 2. Open Docker Desktop manually from Applications, wait for it to start

# 3. Start PostgreSQL
cd /Users/sapaharamsingh/intern_demo/order_management_system
docker compose up -d postgres

# 4. Setup backend
cd backend
cp .env.sample .env
go mod download
go run cmd/main.go --migrate
go run cmd/main.go --api --port=8080

# 5. In a NEW terminal, setup frontend
cd /Users/sapaharamsingh/intern_demo/order_management_system/frontend
cp .env.sample .env
npm install
npm run dev
```

## Verify Everything is Working

1. **Database**: `docker ps | grep postgres` should show a running container
2. **Backend**: `curl http://localhost:8080/api/v1/health` should return "OK"
3. **Frontend**: Open `http://localhost:3000` in your browser


# Quick Install Guide

## Automated Installation

Run the setup script:

```bash
./INSTALL_AND_SETUP.sh
```

This will:
1. Install Docker Desktop (if not already installed)
2. Start PostgreSQL container
3. Create environment configuration files

**Note:** You may need to enter your password for Homebrew installations.

## Manual Installation

If the script doesn't work, follow these steps:

### 1. Install Docker Desktop

```bash
brew install --cask docker
```

Then open Docker Desktop from Applications folder and wait for it to start.

### 2. Start PostgreSQL

```bash
docker compose up -d postgres
```

Or if that doesn't work:
```bash
docker-compose up -d postgres
```

### 3. Verify PostgreSQL is Running

```bash
docker ps | grep postgres
```

### 4. Setup Backend

```bash
cd backend
go mod download
cp .env.sample .env
go run cmd/main.go --migrate
go run cmd/main.go --api --port=8080
```

### 5. Setup Frontend (in new terminal)

```bash
cd frontend
npm install
cp .env.sample .env
npm run dev
```

## Troubleshooting

**Docker not starting:**
- Open Docker Desktop manually from Applications
- Wait for it to fully start (whale icon in menu bar)
- Try `docker ps` to verify it's working

**Permission errors:**
- You may need to run: `sudo chown -R $(whoami) /opt/homebrew/Cellar`
- Or install Docker Desktop manually from: https://www.docker.com/products/docker-desktop/

**Port already in use:**
- Check what's using port 5432: `lsof -i :5432`
- Stop the conflicting service or change the port in `docker-compose.yml`


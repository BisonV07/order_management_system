# How to Check Docker Connection

## Quick Checks

### 1. Check if Docker is Installed
```bash
docker --version
```
✅ **Expected**: `Docker version X.X.X` (you have 29.1.3)

### 2. Check if Docker Daemon is Running
```bash
docker ps
```
✅ **Success**: Shows list of running containers (may be empty)
❌ **Error**: `Cannot connect to the Docker daemon` or `permission denied` → Docker Desktop is not running

### 3. Check Docker Info
```bash
docker info
```
✅ **Success**: Shows Docker system information
❌ **Error**: Cannot connect → Docker Desktop is not running

### 4. Check Docker Compose Services
```bash
# Try both commands (newer Docker uses 'compose' without hyphen)
docker-compose ps
# OR
docker compose ps
```
✅ **Success**: Shows running services
❌ **Empty/Error**: No services running or Docker Compose not found

## Common Issues and Solutions

### Issue 1: "Permission denied" or "Cannot connect to Docker daemon"

**Solution**: Start Docker Desktop
- **macOS**: Open Docker Desktop application from Applications
- **Windows**: Start Docker Desktop from Start Menu
- **Linux**: Start Docker service: `sudo systemctl start docker`

**Verify it's running**:
```bash
docker ps
```

### Issue 2: Docker Desktop is Running but Still Getting Permission Errors

**Solution**: Add your user to docker group (Linux) or restart Docker Desktop

**macOS/Windows**: Restart Docker Desktop

**Linux**:
```bash
sudo usermod -aG docker $USER
# Then log out and log back in, or run:
newgrp docker
```

### Issue 3: Check if PostgreSQL Container is Running

```bash
# List all containers (including stopped)
docker ps -a

# Check for PostgreSQL container
docker ps -a | grep postgres

# Check Docker Compose services
cd /Users/sapaharamsingh/intern_demo/order_management_system
docker-compose ps
# OR
docker compose ps
```

### Issue 4: Test PostgreSQL Connection

If PostgreSQL container is running, test connection:
```bash
# Get container name/ID
docker ps | grep postgres

# Connect to PostgreSQL
docker exec -it <container_name> psql -U postgres

# Or test from host (if port is exposed)
psql -h localhost -p 5432 -U postgres
```

## Step-by-Step Checklist

1. ✅ **Docker Installed**: `docker --version` works
2. ⚠️ **Docker Running**: `docker ps` works (currently failing - need to start Docker Desktop)
3. ⚠️ **Docker Compose**: `docker compose ps` works
4. ⚠️ **PostgreSQL Container**: Check if running with `docker ps`
5. ⚠️ **PostgreSQL Connection**: Test connection to database

## Quick Fix for Your Current Issue

Based on the error you're seeing:
```
permission denied while trying to connect to the docker API at unix:///Users/sapaharamsingh/.docker/run/docker.sock
```

**Action Required**:
1. Open Docker Desktop application on your Mac
2. Wait for it to fully start (whale icon in menu bar should be steady)
3. Run `docker ps` again to verify connection

## Verify Docker Compose Setup

If you have a `docker-compose.yml` file:
```bash
cd /Users/sapaharamsingh/intern_demo/order_management_system
docker-compose up -d
# OR
docker compose up -d
```

Then check status:
```bash
docker-compose ps
# OR
docker compose ps
```


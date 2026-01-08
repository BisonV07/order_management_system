# How to Run the Order Management System

## Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- Docker and Docker Compose (for PostgreSQL)
- PostgreSQL client (optional, for direct database access)

## Step 1: Start PostgreSQL Database

From the project root directory:

**Option A: Using Docker Compose (newer Docker versions)**
```bash
docker compose up -d postgres
```

**Option B: Using docker-compose (legacy)**
```bash
docker-compose up -d postgres
```

**Note:** If you get "command not found", try `docker compose` (without hyphen). Modern Docker Desktop uses `docker compose` as a plugin.

This will start a PostgreSQL container on port 5432.

Verify it's running:
```bash
docker ps
```

**Alternative: If Docker is not available**, you can use a local PostgreSQL installation:
1. Install PostgreSQL locally
2. Create a database: `createdb oms_db`
3. Update `backend/.env` with your local PostgreSQL connection details

## Step 2: Setup Backend

### 2.1 Navigate to backend directory
```bash
cd backend
```

### 2.2 Install Go dependencies
```bash
go mod download
go mod tidy
```

### 2.3 Configure environment
```bash
cp .env.sample .env
```

Edit `.env` file with your database credentials (defaults should work with docker-compose):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=oms_db
DB_SSLMODE=disable
SERVER_PORT=8080
JWT_SECRET=your-secret-key-here
JWT_EXPIRY=24h
LOG_LEVEL=info
```

### 2.4 Run database migrations
```bash
go run cmd/main.go --migrate
```

Or using Make:
```bash
make migrate
```

### 2.5 Start the API server
```bash
go run cmd/main.go --api --port=8080
```

Or using Make:
```bash
make run
```

The API will be available at `http://localhost:8080`

## Step 3: Setup Frontend

### 3.1 Navigate to frontend directory (in a new terminal)
```bash
cd frontend
```

### 3.2 Install dependencies
```bash
npm install
```

### 3.3 Configure environment
```bash
cp .env.sample .env
```

The `.env` file should contain:
```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 3.4 Start the development server
```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## Quick Start (All Commands)

```bash
# Terminal 1: Start database
docker compose up -d postgres
# OR if that doesn't work:
# docker-compose up -d postgres

# Terminal 2: Backend
cd backend
go mod download
cp .env.sample .env
go run cmd/main.go --migrate
go run cmd/main.go --api --port=8080

# Terminal 3: Frontend
cd frontend
npm install
cp .env.sample .env
npm run dev
```

## Verify Setup

1. **Database**: Check if PostgreSQL is running
   ```bash
   docker ps | grep postgres
   ```

2. **Backend API**: Test health endpoint
   ```bash
   curl http://localhost:8080/api/v1/health
   ```
   Should return: `OK`

3. **Frontend**: Open browser to `http://localhost:3000`

## Troubleshooting

### Database Connection Issues
- Ensure PostgreSQL container is running: `docker ps`
- Check database credentials in `backend/.env`
- Verify port 5432 is not in use: `lsof -i :5432`

### Go Module Issues
- Run `go mod tidy` to fix import issues
- Ensure you're in the `backend` directory when running Go commands

### Frontend Build Issues
- Delete `node_modules` and `package-lock.json`, then run `npm install` again
- Check Node.js version: `node --version` (should be 18+)

### Port Conflicts
- Backend default port: 8080 (change in `backend/.env` or use `--port` flag)
- Frontend default port: 3000 (change in `frontend/vite.config.ts`)
- PostgreSQL default port: 5432 (change in `docker-compose.yml`)

## Development Commands

### Backend
```bash
cd backend

# Run server
make run
# or
go run cmd/main.go --api --port=8080

# Run migrations
make migrate
# or
go run cmd/main.go --migrate

# Run tests
make test
# or
go test ./...

# Build binary
make build
# or
go build -o bin/oms cmd/main.go
```

### Frontend
```bash
cd frontend

# Development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

## Stopping Services

```bash
# Stop frontend: Ctrl+C in frontend terminal

# Stop backend: Ctrl+C in backend terminal

# Stop database
docker compose down
# OR
# docker-compose down
# or to remove volumes too
docker compose down -v
```


#!/bin/bash

# Quick Setup Script for OMS
# Run this after Docker is installed

set -e

echo "🚀 Setting up Order Management System..."

# Create server .env file
echo "📝 Creating server/.env..."
cat > server/.env << 'EOF'
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=oms_db
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080

# JWT Configuration
JWT_SECRET=your-secret-key-here-change-this
JWT_EXPIRY=24h

# Logging
LOG_LEVEL=info
EOF

# Create client .env file
echo "📝 Creating client/.env..."
cat > client/.env << 'EOF'
VITE_API_BASE_URL=http://localhost:8080/api/v1
EOF

echo "✅ Configuration files created!"
echo ""
echo "Next steps:"
echo "1. Make sure Docker Desktop is running"
echo "2. Run: docker compose up -d postgres"
echo "3. Start server: cd server && go run cmd/main.go --migrate && go run cmd/main.go --api --port=8080"
echo "4. Start client: cd client && npm install && npm run dev"


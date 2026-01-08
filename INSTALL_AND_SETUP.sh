#!/bin/bash

# Installation and Setup Script for Order Management System
# This script will install Docker, PostgreSQL, and set up the database

set -e

echo "🚀 Starting OMS Setup..."

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "❌ Homebrew is not installed. Please install it first:"
    echo "   /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
    exit 1
fi

echo "✅ Homebrew found"

# Install Docker Desktop
echo ""
echo "📦 Installing Docker Desktop..."
if ! command -v docker &> /dev/null; then
    brew install --cask docker
    echo "✅ Docker Desktop installed"
    echo "⚠️  Please open Docker Desktop from Applications to start it"
    echo "   Waiting 30 seconds for Docker to start..."
    sleep 30
else
    echo "✅ Docker already installed"
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo "⚠️  Docker is not running. Please start Docker Desktop and run this script again."
    exit 1
fi

echo "✅ Docker is running"

# Install PostgreSQL (via Docker, so we don't need to install it separately)
echo ""
echo "✅ PostgreSQL will be run via Docker"

# Start PostgreSQL container
echo ""
echo "🐘 Starting PostgreSQL container..."
cd "$(dirname "$0")"
docker compose up -d postgres

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Check if PostgreSQL is ready
for i in {1..30}; do
    if docker compose exec -T postgres pg_isready -U postgres &> /dev/null; then
        echo "✅ PostgreSQL is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ PostgreSQL failed to start"
        exit 1
    fi
    sleep 1
done

# Setup server environment
echo ""
echo "⚙️  Setting up server configuration..."
cd server
if [ ! -f .env ]; then
    cp .env.sample .env
    echo "✅ Created server/.env file"
else
    echo "✅ server/.env already exists"
fi

# Setup client environment
echo ""
echo "⚙️  Setting up client configuration..."
cd ../client
if [ ! -f .env ]; then
    cp .env.sample .env
    echo "✅ Created client/.env file"
else
    echo "✅ client/.env already exists"
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Server: cd server && go mod download && go run cmd/main.go --migrate && go run cmd/main.go --api --port=8080"
echo "2. Client: cd client && npm install && npm run dev"
echo ""
echo "Database is running at: localhost:5432"
echo "  User: postgres"
echo "  Password: postgres"
echo "  Database: oms_db"


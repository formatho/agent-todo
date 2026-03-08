#!/bin/bash

echo "🔄 Rebuilding Agent Todo Platform (no cache)..."
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Stop and remove containers
echo "🛑 Stopping existing containers..."
docker compose down

# Remove old images
echo "🗑️  Removing old images..."
docker rmi agent-todo-backend agent-todo-frontend 2>/dev/null || true

# Build images without cache
echo "🔨 Building new images (no cache)..."
docker compose build --no-cache

# Start services
echo "🚀 Starting services..."
docker compose up -d

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check if services are healthy
echo ""
echo "📊 Service Status:"
docker compose ps

echo ""
echo "✅ Rebuild complete!"
echo ""
echo "📍 Access Points:"
echo "   Frontend:  http://localhost:3000"
echo "   Backend:   http://localhost:8080"
echo "   Swagger:   http://localhost:8080/docs"
echo ""

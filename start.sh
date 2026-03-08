#!/bin/bash

echo "🚀 Starting Agent Todo Platform..."
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Start services
docker compose up -d

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check if services are healthy
echo ""
echo "📊 Service Status:"
docker compose ps

echo ""
echo "✅ Agent Todo Platform is starting!"
echo ""
echo "📍 Access Points:"
echo "   Frontend:  http://localhost:3000"
echo "   Backend:   http://localhost:8080"
echo "   Swagger:   http://localhost:8080/docs"
echo "   Health:    http://localhost:8080/health"
echo ""
echo "🔑 Default Credentials:"
echo "   Email:    admin@example.com"
echo "   Password: admin123"
echo ""
echo "📝 Useful Commands:"
echo "   View logs:  docker compose logs -f"
echo "   Stop:       ./stop.sh or docker compose down"
echo ""

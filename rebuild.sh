#!/bin/bash

echo "🔄 Rebuilding Agent Todo Platform with new UI..."
echo ""

# Stop running containers
echo "🛑 Stopping containers..."
docker compose down

# Remove old images (force rebuild)
echo "🗑️  Removing old images..."
docker compose build --no-cache

# Start services
echo "🚀 Starting services..."
docker compose up

# If user presses Ctrl+C, stop containers
trap 'echo ""; echo "🛑 Stopping services..."; docker compose down; exit 0' INT

#!/bin/bash

echo "🚀 Starting Agent Todo Platform..."
echo ""

# Check if docker compose is available
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Stop any running containers
echo "🛑 Stopping any running containers..."
docker compose down 2>/dev/null

# Build and start services
echo "🔨 Building and starting services..."
echo ""
docker compose up --build

# If user presses Ctrl+C, stop containers
trap 'echo ""; echo "🛑 Stopping services..."; docker compose down; exit 0' INT

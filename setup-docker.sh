#!/bin/bash

# LunchRush Docker Setup Script

echo "🍱 LunchRush Docker Setup"
echo "========================="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker not found. Please install Docker first:"
    echo "   https://www.docker.com/get-started"
    exit 1
fi

echo "✅ Docker installed: $(docker --version)"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

echo "✅ Docker is running"

# Check if docker-compose is available
if command -v docker-compose &> /dev/null; then
    echo "✅ Docker Compose installed: $(docker-compose --version)"
elif docker compose version &> /dev/null; then
    echo "✅ Docker Compose (plugin) installed: $(docker compose version)"
    # Create alias for docker-compose
    alias docker-compose='docker compose'
else
    echo "❌ Docker Compose not found. Please install Docker Compose."
    exit 1
fi

# Make scripts executable
chmod +x start-docker.sh stop-docker.sh setup-docker.sh

echo ""
echo "✅ Docker setup complete!"
echo ""
echo "🚀 To start LunchRush with Docker:"
echo "   ./start-docker.sh"
echo ""
echo "🛑 To stop LunchRush:"
echo "   ./stop-docker.sh"
echo ""
echo "📋 Manual commands:"
echo "   docker-compose up --build    → Build and start"
echo "   docker-compose down          → Stop services"
echo "   docker-compose logs -f       → View logs"
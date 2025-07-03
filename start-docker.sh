#!/bin/bash

# LunchRush Docker Startup Script

echo "🍱 Starting LunchRush with Docker..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Stop and remove existing containers if any
echo "🧹 Cleaning up existing containers..."
docker-compose down --remove-orphans

# Build and start all services
echo "🚀 Building and starting all services..."
docker-compose up --build -d

# Wait for services to start
echo "⏳ Waiting for services to start..."
sleep 10

# Check if services are running
echo "📊 Checking service status..."
docker-compose ps

# Check backend health
echo "🔍 Checking backend health..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Backend is healthy!"
else
    echo "⚠️  Backend might still be starting up..."
fi

echo ""
echo "✅ LunchRush is running in background!"
echo ""
echo "📍 Frontend: http://localhost:8000"
echo "📍 Backend API: http://localhost:8080"
echo "📍 Redis: localhost:6379"
echo ""
echo "📋 Useful commands:"
echo "   docker-compose logs -f           → Follow logs"
echo "   docker-compose logs lunchservice → Backend logs"
echo "   docker-compose ps                → Check status"
echo "   docker-compose down              → Stop all services"
echo ""
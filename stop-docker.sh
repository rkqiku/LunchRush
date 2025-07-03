#!/bin/bash

# LunchRush Docker Stop Script

echo "🛑 Stopping LunchRush Docker services..."

# Stop all services
docker-compose down

# Optional: Remove volumes
read -p "Remove Redis data volume? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker-compose down -v
    echo "✅ All services and volumes stopped"
else
    echo "✅ All services stopped (data preserved)"
fi

# Show remaining containers (if any)
echo ""
echo "📊 Remaining LunchRush containers:"
docker ps -a | grep -E "(lunchservice|redis|frontend)" || echo "None"
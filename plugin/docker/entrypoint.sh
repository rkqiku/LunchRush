#!/bin/bash
set -e

# Function to handle signals
shutdown() {
    echo "Shutting down nginx..."
    nginx -s quit
    exit 0
}

# Trap signals
trap shutdown SIGTERM SIGINT

# Execute all scripts in /docker-entrypoint.d/
if [ -d /docker-entrypoint.d ]; then
    for f in /docker-entrypoint.d/*.sh; do
        if [ -x "$f" ]; then
            echo "Running $f..."
            "$f"
        fi
    done
fi

# Process nginx config templates with environment variables
# Note: The default nginx entrypoint already handles templates, but we'll ensure
# our variables are set with defaults
export NGINX_PORT="${NGINX_PORT:-80}"
export NGINX_HOST="${NGINX_HOST:-localhost}"
export API_PROXY_URL="${API_PROXY_URL:-http://localhost:8080}"

# Start nginx
echo "Starting nginx..."
exec "$@"
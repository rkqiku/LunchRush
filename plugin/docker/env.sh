#!/bin/bash
# Script to inject runtime environment variables into the React app

# Default values
: ${REACT_APP_API_URL:=http://localhost:8080}
: ${REACT_APP_ENVIRONMENT:=production}

# Create runtime config file
cat > /usr/share/nginx/html/env-config.js <<EOF
window._env_ = {
  REACT_APP_API_URL: "${REACT_APP_API_URL}",
  REACT_APP_ENVIRONMENT: "${REACT_APP_ENVIRONMENT}",
  REACT_APP_VERSION: "${REACT_APP_VERSION:-1.0.0}",
  REACT_APP_FEATURES: {
    ENABLE_VOTING: ${REACT_APP_ENABLE_VOTING:-true},
    ENABLE_NOTIFICATIONS: ${REACT_APP_ENABLE_NOTIFICATIONS:-true}
  }
};
EOF

echo "Runtime environment variables injected successfully:"
echo "- REACT_APP_API_URL: ${REACT_APP_API_URL}"
echo "- REACT_APP_ENVIRONMENT: ${REACT_APP_ENVIRONMENT}"
# syntax=docker/dockerfile:1.4

# Build stage - Dependencies
FROM node:20-alpine AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Copy dependency files
COPY package*.json ./

# Install dependencies with clean install
RUN npm ci --only=production && \
    cp -R node_modules prod_node_modules && \
    npm ci

# Build stage - Builder
FROM node:20-alpine AS builder
WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules
COPY package*.json ./

# Copy source code
COPY . .

# Build arguments for compile-time variables
ARG VITE_API_URL=http://localhost:8080
ARG NODE_ENV=production

# Set environment variables for build
ENV VITE_API_URL=$VITE_API_URL
ENV NODE_ENV=$NODE_ENV

# Build the application
RUN npm run build

# Runtime stage - nginx
FROM nginx:stable-alpine AS runtime

# Install runtime dependencies
RUN apk add --no-cache --update \
    bash \
    curl \
    ca-certificates \
    gettext \
    && rm -rf /var/cache/apk/*

# Copy nginx configuration
COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY nginx/default.conf.template /etc/nginx/templates/default.conf.template

# Copy built application
COPY --from=builder --chown=nginx:nginx /app/dist /usr/share/nginx/html

# Copy runtime script
COPY docker/env.sh /docker-entrypoint.d/01-env.sh
RUN chmod +x /docker-entrypoint.d/01-env.sh

# Create env-config.js placeholder
RUN touch /usr/share/nginx/html/env-config.js && \
    chown nginx:nginx /usr/share/nginx/html/env-config.js

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/health || exit 1

# Expose port (informational)
EXPOSE 80

# Let nginx handle the entrypoint
CMD ["nginx", "-g", "daemon off;"]
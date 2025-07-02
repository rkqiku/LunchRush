# 🍱 LunchRush - Collaborative Lunch Ordering System

> **🧠 Important Note:** I am not a software engineer with experience in React or Go. This entire project was built with AI assistance in approximately 4 hours, demonstrating that with the right tools, anyone can create a working MVP in record time. The barriers to entry in software development are disappearing - this project is proof that determination and AI can bridge any knowledge gap.

A real-time collaborative lunch ordering system with restaurant voting, meal tracking, and automatic session management. Built with modern microservices architecture using Go, React, Docker, and Dapr.

## 🚀 Quick Start

### Prerequisites
- **Docker** (that's it! No other dependencies needed)

### Installation & Running

```bash
# 1. Clone repository
git clone https://github.com/hysaordis/LunchRush.git
cd LunchRush

# 2. Setup (first time only - checks Docker installation)
./setup-docker.sh

# 3. Start all services
./start-docker.sh

# 4. Access the application
open http://localhost:8000
```

### Stopping Services

```bash
# Stop all services (preserves data)
./stop-docker.sh

# Stop and remove all data
./stop-docker.sh  # then answer 'y' when prompted
```

## 🎯 Key Features

- **Real-time Collaboration** - See who's joining and ordering in real-time
- **Restaurant Voting** - Propose and vote for restaurants democratically
- **Order Management** - Track individual meal choices with live updates
- **Session Locking** - Automatic or manual session locking at specified times
- **User Activity Tracking** - Heartbeat system shows active/inactive users
- **Order Assignment** - Designate who will place the group order
- **Automatic Cleanup** - Inactive users are marked and can be removed

## 🏗️ Architecture

### Services Overview

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend  │────▶│   Backend   │────▶│    Redis    │
│  (React)    │     │    (Go)     │     │   (State)   │
│  Port 8000  │     │  Port 8080  │     │  Port 6379  │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                    ┌──────▼──────┐
                    │    Dapr     │
                    │  (Sidecar)  │
                    └─────────────┘
```

### Tech Stack

- **Frontend**: React + Vite + Tailwind CSS + React Query
- **Backend**: Go + Chi Router + Dapr SDK
- **State Management**: Redis via Dapr State Store
- **Real-time Updates**: Polling (2s interval) + Pub/Sub events
- **Infrastructure**: Docker Compose with health checks

### Docker Implementation

- **Multi-stage builds** for optimized images
- **Runtime environment injection** for flexible deployments
- **Security headers** (CSP, HSTS, X-Frame-Options, etc.)
- **Health checks** for all services
- **Automatic restart** policies
- **Resource limits** and logging configuration

## 📡 API Documentation

### Session Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/sessions` | Create a new lunch session |
| GET | `/sessions/today` | Get today's session |
| GET | `/sessions/{id}` | Get specific session details |
| POST | `/sessions/{id}/join` | Join a session |
| PUT | `/sessions/{id}/meal` | Update meal choice |
| POST | `/sessions/{id}/lock` | Lock the session |
| POST | `/sessions/{id}/heartbeat` | Update user activity status |
| DELETE | `/sessions/{id}/participants/{username}` | Remove a participant |

### Restaurant Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/sessions/{id}/restaurants` | Get all proposed restaurants |
| POST | `/sessions/{id}/restaurants` | Propose a new restaurant |
| POST | `/sessions/{id}/restaurants/{restaurantId}/vote` | Vote/unvote for a restaurant |
| DELETE | `/sessions/{id}/restaurants/{restaurantId}` | Delete a restaurant proposal |

### Order Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/sessions/{id}/order-placer` | Assign who will place the order |

### Example Requests

```bash
# Create a session
curl -X POST http://localhost:8080/sessions

# Join session
curl -X POST http://localhost:8080/sessions/{session-id}/join \
  -H "Content-Type: application/json" \
  -d '{"username": "John Doe"}'

# Propose restaurant
curl -X POST http://localhost:8080/sessions/{session-id}/restaurants \
  -H "Content-Type: application/json" \
  -d '{"name": "Pizza Palace", "proposedBy": "John Doe"}'

# Vote for restaurant
curl -X POST http://localhost:8080/sessions/{session-id}/restaurants/{restaurant-id}/vote \
  -H "Content-Type: application/json" \
  -d '{"username": "John Doe"}'

# Update meal choice
curl -X PUT http://localhost:8080/sessions/{session-id}/meal \
  -H "Content-Type: application/json" \
  -d '{"username": "John Doe", "meal": "Margherita Pizza"}'
```

### Pub/Sub Events

The backend publishes the following events via Dapr:

- `session.created` - New session created
- `session.locked` - Session locked for ordering
- `participant.joined` - User joined session
- `participant.left` - User removed from session
- `meal.updated` - Meal choice updated
- `restaurant.proposed` - New restaurant suggested
- `restaurant.voted` - Vote cast for restaurant
- `orderplacer.set` - Order placer assigned

## 🔧 Development

### Local Development

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Rebuild after code changes
docker-compose up --build

# Check service health
curl http://localhost:8080/health
curl http://localhost:8000/health
```

### Backend Development (without Docker)

Requires Go 1.21+, Dapr CLI, and Redis:

```bash
cd microservice
make deps
make run-dapr
```

### Frontend Development (without Docker)

Requires Node.js 20+:

```bash
cd plugin
npm install
npm run dev
```

### Environment Variables

Configure via `.env` file or environment:

```bash
# Backend API URL (for frontend)
REACT_APP_API_URL=http://localhost:8080

# Environment name
REACT_APP_ENVIRONMENT=development

# Feature flags
REACT_APP_ENABLE_VOTING=true
REACT_APP_ENABLE_NOTIFICATIONS=true
```

### Building for Production

```bash
# Build with custom configuration
docker-compose build \
  --build-arg VITE_API_URL=https://api.production.com \
  frontend

# Run with runtime configuration
REACT_APP_API_URL=https://api.production.com \
docker-compose up -d
```

## 🧪 Testing

### Manual Testing

1. Start services: `./start-docker.sh`
2. Open multiple browser windows at http://localhost:8000
3. Create a session in one window
4. Join from other windows using different usernames
5. Test features:
   - Propose and vote for restaurants
   - Add meal choices
   - Watch real-time updates
   - Test user inactivity (wait 5 minutes)
   - Lock session and verify restrictions

### API Testing

Use the provided curl examples or import into Postman/Insomnia.

## 📊 Monitoring

- **Health Endpoints**: 
  - Backend: http://localhost:8080/health
  - Frontend: http://localhost:8000/health
- **Logs**: Structured JSON logging with rotation
- **Metrics**: Docker stats and health check status

## 🛡️ Security Features

- Non-root container execution
- Comprehensive security headers
- Rate limiting on API endpoints
- Input validation and sanitization
- CORS properly configured
- No hardcoded secrets

## 📋 Manual Docker Commands

If you prefer not to use the convenience scripts:

```bash
# Start all services
docker-compose up --build -d

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# View logs
docker-compose logs -f

# Restart a specific service
docker-compose restart backend

# Check service status
docker-compose ps
```

## 🚀 Deployment

The application is designed to be deployed anywhere Docker runs:

1. **Local Development**: Use docker-compose as described
2. **Cloud Deployment**: Deploy to any container service (ECS, GKE, AKS)
3. **Kubernetes**: Convert docker-compose to K8s manifests
4. **Single Server**: Run docker-compose in production mode

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📝 License

[Add your license here]

---

**Built for the Huly Plugin Challenge** 🎯

*Developed with AI assistance, proving that technology barriers are meant to be broken! 🚀*